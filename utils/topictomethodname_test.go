package utils

import (
	"testing"

	"github.com/0xanonymeow/go-subtest"
)

func TestTopicToMethodName(t *testing.T) {
	subtests := []subtest.Subtest{
		{
			Name:         "topic_to_method_name",
			ExpectedData: "TestExampleHandler",
			ExpectedErr:  nil,
			Test: func() (interface{}, error) {
				return TopicToMethodName("dev.service-test.test-example"), nil
			},
		},
		{
			Name:         "topic_to_method_name_no_dot",
			ExpectedData: "TopicWithoutDotHandler",
			ExpectedErr:  nil,
			Test: func() (interface{}, error) {
				return TopicToMethodName("topic-without-dot"), nil
			},
		},
	}

	subtest.RunSubtests(t, subtests)
}
