package configuration_test

import (
    "strings"
    "testing"

    "github.com/form3tech-oss/github-team-approver-commons/pkg/configuration"
    "github.com/magiconair/properties/assert"
    "github.com/stretchr/testify/require"
)

func Test_ReadConfiguration(t *testing.T) {
    tt := []struct {
        name     string
        payload  string
        expected configuration.Configuration
    }{
        {
            name: "basic",
            payload: `
ignore_contributor_approval: false
pull_request_approval_rules:
    - rules:
       - regex: "regex"
         approving_team_handles:
         - team-a
         - team-b`,
            expected: configuration.Configuration{
                IgnoreContributorApproval: false,
                PullRequestApprovalRules: []configuration.PullRequestApprovalRule{
                    {
                        Rules: []configuration.Rule{
                            {
                                Regex: "regex",
                                ApprovingTeamHandles: []string{
                                    "team-a",
                                    "team-b",
                                },
                            },
                        },
                    },
                },
            },
        },
    }

    for _, tc := range tt {
        t.Run(tc.name, func(t *testing.T) {
            strings.NewReader(tc.payload)
            cfg, err := configuration.ReadConfiguration(strings.NewReader(tc.payload))
            require.NoError(t, err)
            require.NotNil(t, cfg)
            assert.Equal(t, *cfg, tc.expected)
        })
    }
}
