package configuration_test

import (
    "os"
    "testing"

    "github.com/form3tech-oss/github-team-approver-commons/pkg/configuration"
    "github.com/magiconair/properties/assert"
    "github.com/stretchr/testify/require"
)

func Test_ReadConfiguration(t *testing.T) {
    tt := []struct {
        name        string
        payloadPath string
        expected    configuration.Configuration
    }{
        {
            name:        "basic",
            payloadPath: "./testdata/configuration_basic.yaml",
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
            cfgFile, err := os.Open(tc.payloadPath)
            cfg, err := configuration.ReadConfiguration(cfgFile)
            require.NoError(t, err)
            require.NotNil(t, cfg)
            assert.Equal(t, *cfg, tc.expected)
        })
    }
}
