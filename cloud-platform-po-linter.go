package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/ghodss/yaml"
	"github.com/ministryofjustice/cloud-platform-environments/pkg/authenticate"
	"github.com/prometheus-operator/prometheus-operator/pkg/admission"
	v1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	"github.com/prometheus/prometheus/model/rulefmt"
	githubaction "github.com/sethvargo/go-githubactions"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func prDir(token, githubrepo, githubref string) ([]string, error) {
	if token == "" {
		return nil, errors.New("you must have a valid GitHub token")
	}
	if githubrepo == "" {
		return nil, errors.New("you must have a valid GitHub Repository")
	}
	if githubref == "" {
		return nil, errors.New("you must have a valid GitHub Ref")
	}

	client, err := authenticate.GitHubClient(token)
	if err != nil {
		return nil, err
	}

	//repo user and repo name
	githubrepoS := strings.Split(githubrepo, "/")
	owner := githubrepoS[0]
	repo := githubrepoS[1]

	// get pr owner
	githubrefS := strings.Split(githubref, "/")
	branch := githubrefS[2]
	bid, err := strconv.Atoi(branch)
	if err != nil {
		log.Fatal(err)
	}

	repos, _, _ := client.PullRequests.ListFiles(context.Background(), owner, repo, bid, nil)

	var fileDir []string
	for _, repo := range repos {
		if strings.Contains(*repo.Filename, directory) {
			fileDir = append(fileDir, *repo.Filename)
		}
	}
	return fileDir, nil

}

func poLint(dir []string) error {
	files := dir

	for _, filename := range files {
		log.SetPrefix(fmt.Sprintf("%s: ", filename))
		content, err := os.ReadFile(filename)
		if err != nil {
			log.Fatal(err)
		}

		var meta metav1.TypeMeta

		err = yaml.Unmarshal(content, &meta)
		if err != nil {
			return fmt.Errorf("%v", err)
		}

		switch meta.Kind {
		case v1.PrometheusRuleKind:
			j, err := yaml.YAMLToJSON(content)
			if err != nil {
				fmt.Println("unable to convert YAML to JSON: %w", err)
				return fmt.Errorf("unable to convert YAML to JSON: %w", err)
			}

			decoder := json.NewDecoder(bytes.NewBuffer(j))
			decoder.DisallowUnknownFields()

			var rule v1.PrometheusRule
			err = decoder.Decode(&rule)
			if err != nil {
				fmt.Println("prometheus rule is invalid: %w", err)
				return fmt.Errorf("prometheus rule is invalid: %w", err)
			}
			err = validateRules(content)
			if err != nil {
				fmt.Println("prometheus rule validation failed: %w", err)
				return fmt.Errorf("prometheus rule validation failed: %w", err)
			}
		default:
			return fmt.Errorf("MetaType is unknown to linter. Not in PrometheusRule")
		}
	}
	return nil
}

func validateRules(content []byte) error {
	rule := &admission.PrometheusRules{}
	err := yaml.Unmarshal(content, rule)
	if err != nil {
		return fmt.Errorf("unable load prometheus rule: %w", err)
	}
	rules, errorsArray := rulefmt.Parse(rule.Spec.Raw)
	if len(errorsArray) != 0 {
		for _, err := range errorsArray {
			log.Println(err)
			return fmt.Errorf("%w", err)
		}
	}
	if len(rules.Groups) == 0 {
		return errors.New("no group found")
	}
	for _, group := range rules.Groups {
		if len(group.Rules) == 0 {
			return fmt.Errorf("no rules found in group: %s: %w", group.Name, err)
		}
	}
	return nil
}

var (
	token      = flag.String("token", os.Getenv("GITHUB_OAUTH_TOKEN"), "GihHub Personel token string")
	githubrepo = flag.String("githubrepo", os.Getenv("GITHUB_REPOSITORY"), "Github Repository string")
	githubref  = flag.String("githubref", os.Getenv("GITHUB_REF"), "Github Respository PR ref string")
	directory  = os.Getenv("REPO_DIRECTORY")
)

func main() {
	d, err := prDir(*token, *githubrepo, *githubref)
	if err != nil {
		log.Fatal(err)
	}
	err = poLint(d)
	if err != nil {
		o := fmt.Sprintf("%v", err)
		githubaction.SetOutput("po-linter", o)
		os.Exit(1)
	}
}
