package configuration

import (
	"fmt"
	"io"

	"gopkg.in/yaml.v2"
)

const (
	// ConfigurationFilePath is the path to the file where configuration is specified, relative to the root of the repository.
	ConfigurationFilePath = ".github/GITHUB_TEAM_APPROVER.yaml"
)

type ApprovalMode string

func (am *ApprovalMode) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var v string
	if err := unmarshal(&v); err != nil {
		return err
	}
	m := ApprovalMode(v)
	switch m {
	case ApprovalModeRequireAll, ApprovalModeRequireAny:
		*am = m
		return nil
	default:
		return fmt.Errorf("%q is not a valid approval mode", v)
	}
}

const (
	// ApprovalModeRequireAll indicates that a review from all of the listed teams is required for a PR to be approved.
	ApprovalModeRequireAll ApprovalMode = "require_all"
	// ApprovalModeRequireAny indicates that a review from any of the listed teams is required for a PR to be approved.
	ApprovalModeRequireAny ApprovalMode = "require_any"
)

// Configuration is used to configure the approval process for PRs in a given repository.
type Configuration struct {
	PullRequestApprovalRules []PullRequestApprovalRule `yaml:"pull_request_approval_rules"`
}

// PullRequestApprovalRule is used to associate a set of rules with a set of target branches.
type PullRequestApprovalRule struct {
	// Rules is the set of rules applied to PRs targeting a specific branch.
	Rules []Rule `yaml:"rules"`
	// Alerts to fire when PR is merged successfully
	Alerts []Alert `yaml:"alerts"`
	// TargetBranch is the target branch of the PR.
	TargetBranches []string `yaml:"target_branches"`
}

// Rule is used to configure approval of PRs based on a particular regular expression.
type Rule struct {
	// Regex is the regular expression to match each PR's body against for the current check.
	Regex string `yaml:"regex,omitempty"`
	// RegexLabel is regular expression to match each PR's labels for the current check
	RegexLabel string `yaml:"regex_label,omitempty"`
	// Directories to check for changes, leave empty to check all directories.
	Directories []string `yaml:"directories,omitempty"`
	// ApprovingTeamHandles is the list of IDs/slugs/names of the teams that must approve each matched PR.
	ApprovingTeamHandles []string `yaml:"approving_team_handles"`
	// ApprovalMode specifies the approval mode for PRs that match this check.
	ApprovalMode ApprovalMode `yaml:"approval_mode"`
	// Labels is the set of labels to add to the PR according to whether the PR's body matches the regular expression above.
	Labels []string `yaml:"labels"`
	// ForceApproval indicates whether to forcibly approve the PR regardless of the current status of reviews.
	ForceApproval bool `yaml:"force_approval"`
	// IgnoreContributorApproval will enforce that a reviewer can not approve a pull request that they have contributed
	// towards. That is, a reviewer's approval is only considered if:
	// - they have _not_ pushed a commit to the branch being merged,
	// - they have _not_ co-authored a commit in the branch being merged.
	// This does not include UI merges from the repositories main branch.
	// See https://github.com/form3tech-oss/github-team-approver/pull/27 for more details.
	// Defaults to false.
	IgnoreContributorApproval bool `yaml:"ignore_contributor_approval"`
}

type Alert struct {
	// Regex is the regular expression to match each PR's body against for the alert
	Regex string `yaml:"regex"`
	// Slack webhook url secret - the location of the secret where the slack webhook url is stored - Slack webhooks contain sensitive data so must be hidden
	SlackWebhookSecret string `yaml:"slack_webhook_secret"`
	// Slack message to send to slack
	SlackMessage string `yaml:"slack_message"`
}

// ReadConfiguration attempts to read a Configuration object from the provided Reader.
func ReadConfiguration(r io.Reader) (*Configuration, error) {
	v := &Configuration{}
	d := yaml.NewDecoder(r)
	if err := d.Decode(v); err != nil {
		return nil, fmt.Errorf("error decoding configuration: %v", err)
	}
	return v, nil
}

// Write attempts to write the current Configuration object to the provided Writer.
func (c *Configuration) Write(w io.Writer) error {
	e := yaml.NewEncoder(w)
	defer e.Close()
	return e.Encode(c)
}
