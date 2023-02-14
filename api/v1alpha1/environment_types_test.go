package v1alpha1

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGitCommitEquals(t *testing.T) {
	testCases := []struct {
		name           string
		lhs            *GitCommit
		rhs            *GitCommit
		expectedResult bool
	}{
		{
			name:           "lhs and rhs both nil",
			expectedResult: true,
		},
		{
			name:           "only lhs is nil",
			rhs:            &GitCommit{},
			expectedResult: false,
		},
		{
			name:           "only rhs is nil",
			lhs:            &GitCommit{},
			expectedResult: false,
		},
		{
			name: "repoUrls differ",
			lhs: &GitCommit{
				RepoURL: "foo",
				ID:      "fake-commit-id",
			},
			rhs: &GitCommit{
				RepoURL: "bar",
				ID:      "fake-commit-id",
			},
			expectedResult: false,
		},
		{
			name: "commit IDs differ",
			lhs: &GitCommit{
				RepoURL: "fake-url",
				ID:      "foo",
			},
			rhs: &GitCommit{
				RepoURL: "fake-url",
				ID:      "bar",
			},
			expectedResult: false,
		},
		{
			name: "perfect match",
			lhs: &GitCommit{
				RepoURL: "fake-url",
				ID:      "fake-commit-id",
			},
			rhs: &GitCommit{
				RepoURL: "fake-url",
				ID:      "fake-commit-id",
			},
			expectedResult: true,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			require.Equal(
				t,
				testCase.expectedResult,
				testCase.lhs.Equals(testCase.rhs),
			)
		})
	}
}

func TestEnvironmentStateSameMaterials(t *testing.T) {
	testCases := []struct {
		name           string
		lhs            *EnvironmentState
		rhs            *EnvironmentState
		expectedResult bool
	}{
		{
			name:           "lhs and rhs both nil",
			expectedResult: true,
		},
		{
			name:           "only lhs is nil",
			rhs:            &EnvironmentState{},
			expectedResult: false,
		},
		{
			name:           "only rhs is nil",
			lhs:            &EnvironmentState{},
			expectedResult: false,
		},
		{
			name: "git commits differ",
			lhs: &EnvironmentState{
				GitCommit: &GitCommit{
					RepoURL: "fake-url",
					ID:      "old-commit",
				},
			},
			rhs: &EnvironmentState{
				GitCommit: &GitCommit{
					RepoURL: "fake-url",
					ID:      "new-commit",
				},
			},
			expectedResult: false,
		},
		{
			name: "images have different cardinality",
			lhs:  &EnvironmentState{},
			rhs: &EnvironmentState{
				Images: []Image{
					{
						RepoURL: "nginx",
						Tag:     "1.23.3",
					},
				},
			},
			expectedResult: false,
		},
		{
			name: "charts have different cardinality",
			lhs:  &EnvironmentState{},
			rhs: &EnvironmentState{
				Charts: []Chart{
					{
						RegistryURL: "fake-registry",
						Name:        "fake-chart",
						Version:     "fake-version",
					},
				},
			},
			expectedResult: false,
		},
		{
			name: "images have same cardinality, but differ",
			lhs: &EnvironmentState{
				Images: []Image{
					{
						RepoURL: "nginx",
						Tag:     "1.23.2",
					},
				},
			},
			rhs: &EnvironmentState{
				Images: []Image{
					{
						RepoURL: "nginx",
						Tag:     "1.23.3",
					},
				},
			},
			expectedResult: false,
		},
		{
			name: "charts have same cardinality, but differ",
			lhs: &EnvironmentState{
				Charts: []Chart{
					{
						RegistryURL: "fake-registry",
						Name:        "fake-chart",
						Version:     "fake-version",
					},
				},
			},
			rhs: &EnvironmentState{
				Charts: []Chart{
					{
						RegistryURL: "fake-registry",
						Name:        "fake-chart",
						Version:     "different-fake-version",
					},
				},
			},
			expectedResult: false,
		},
		{
			name: "perfect match",
			lhs: &EnvironmentState{
				Images: []Image{
					{
						RepoURL: "foo",
						Tag:     "1.0.0",
					},
					{
						RepoURL: "bar",
						Tag:     "1.0.0",
					},
				},
				Charts: []Chart{
					{
						RegistryURL: "fake-registry",
						Name:        "fake-chart",
						Version:     "fake-version",
					},
					{
						RegistryURL: "another-fake-registry",
						Name:        "another-fake-chart",
						Version:     "another-fake-version",
					},
				},
			},
			rhs: &EnvironmentState{
				// Note that we make a point of putting the images in different orders
				// here, because order shouldn't matter.
				Images: []Image{
					{
						RepoURL: "bar",
						Tag:     "1.0.0",
					},
					{
						RepoURL: "foo",
						Tag:     "1.0.0",
					},
				},
				// Note that we make a point of putting the charts in different orders
				// here, because order shouldn't matter.
				Charts: []Chart{
					{
						RegistryURL: "another-fake-registry",
						Name:        "another-fake-chart",
						Version:     "another-fake-version",
					},
					{
						RegistryURL: "fake-registry",
						Name:        "fake-chart",
						Version:     "fake-version",
					},
				},
			},
			expectedResult: true,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			require.Equal(
				t,
				testCase.expectedResult,
				testCase.lhs.SameMaterials(testCase.rhs),
			)
		})
	}
}