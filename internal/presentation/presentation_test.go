package presentation_test

import (
	"bytes"
	"fmt"
	"io"
	"testing"

	"github.com/LLIEPJIOK/mazegenerator/internal/domain"
	"github.com/LLIEPJIOK/mazegenerator/internal/presentation"
	"github.com/stretchr/testify/require"
)

func TestProcessInput(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		input    string
		expected *presentation.Input
	}{
		{
			input: "10\n10\n0\n0\n9\n9\n1\n1",
			expected: presentation.NewInput(
				10,
				10,
				domain.NewCoord(0, 0),
				domain.NewCoord(9, 9),
				"prim",
				"dijkstra",
			),
		},
		{
			input: "15\n15\n2\n0\n14\n14\n1\n2",
			expected: presentation.NewInput(
				15,
				15,
				domain.NewCoord(2, 0),
				domain.NewCoord(14, 14),
				"prim",
				"a-star",
			),
		},
		{
			input: "20\n20\n5\n19\n19\n2\n2\n1",
			expected: presentation.NewInput(
				20,
				20,
				domain.NewCoord(5, 19),
				domain.NewCoord(19, 2),
				"backtrack",
				"dijkstra",
			),
		},
		{
			input: "12\n12\n0\n11\n11\n0\n2\n2",
			expected: presentation.NewInput(
				12,
				12,
				domain.NewCoord(0, 11),
				domain.NewCoord(11, 0),
				"backtrack",
				"a-star",
			),
		},
		{
			input: "30\n30\n0\n15\n29\n18\n1\n1",
			expected: presentation.NewInput(
				30,
				30,
				domain.NewCoord(0, 15),
				domain.NewCoord(29, 18),
				"prim",
				"dijkstra",
			),
		},
		{
			input: "25\n25\n10\n24\n24\n24\n1\n2",
			expected: presentation.NewInput(
				25,
				25,
				domain.NewCoord(10, 24),
				domain.NewCoord(24, 24),
				"prim",
				"a-star",
			),
		},
		{
			input: "18\n18\n3\n0\n0\n17\n2\n1",
			expected: presentation.NewInput(
				18,
				18,
				domain.NewCoord(3, 0),
				domain.NewCoord(0, 17),
				"backtrack",
				"dijkstra",
			),
		},
		{
			input: "8\n8\n0\n0\n7\n7\n2\n2",
			expected: presentation.NewInput(
				8,
				8,
				domain.NewCoord(0, 0),
				domain.NewCoord(7, 7),
				"backtrack",
				"a-star",
			),
		},
		{
			input: "50\n50\n25\n0\n49\n49\n1\n1",
			expected: presentation.NewInput(
				50,
				50,
				domain.NewCoord(25, 0),
				domain.NewCoord(49, 49),
				"prim",
				"dijkstra",
			),
		},
		{
			input: "40\n40\n20\n39\n39\n39\n1\n2",
			expected: presentation.NewInput(
				40,
				40,
				domain.NewCoord(20, 39),
				domain.NewCoord(39, 39),
				"prim",
				"a-star",
			),
		},
		{
			input: "5\n5\n0\n1\n4\n4\n2\n1",
			expected: presentation.NewInput(
				5,
				5,
				domain.NewCoord(0, 1),
				domain.NewCoord(4, 4),
				"backtrack",
				"dijkstra",
			),
		},
	}

	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("#%d", i+1), func(t *testing.T) {
			t.Parallel()

			inputStream := bytes.NewBufferString(testCase.input)
			pres := presentation.New(inputStream, io.Discard)
			got, err := pres.ProcessInput()

			require.NoError(t, err, "input should get without error")
			require.Equal(t, testCase.expected, got, "values should be equal")
		})
	}
}

func TestProcessInputWithInvalidData(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		input    string
		expected *presentation.Input
	}{
		{
			input: "10\n10\n-1\n0\n0\n9\n9\n1\n1",
			expected: presentation.NewInput(
				10,
				10,
				domain.NewCoord(0, 0),
				domain.NewCoord(9, 9),
				"prim",
				"dijkstra",
			),
		},
		{
			input: "15\n15\n55\n2\n0\n14\n14\n1\n2",
			expected: presentation.NewInput(
				15,
				15,
				domain.NewCoord(2, 0),
				domain.NewCoord(14, 14),
				"prim",
				"a-star",
			),
		},
		{
			input: "20\n20\n5\n19\n19\n2\n2\n30\n1",
			expected: presentation.NewInput(
				20,
				20,
				domain.NewCoord(5, 19),
				domain.NewCoord(19, 2),
				"backtrack",
				"dijkstra",
			),
		},
		{
			input: "12\n12\n0\n11\n0\n11\n11\n0\n2\n2",
			expected: presentation.NewInput(
				12,
				12,
				domain.NewCoord(0, 11),
				domain.NewCoord(11, 0),
				"backtrack",
				"a-star",
			),
		},
		{
			input: "30\n30\n4\n4\n0\n15\n29\n18\n1\n1",
			expected: presentation.NewInput(
				30,
				30,
				domain.NewCoord(0, 15),
				domain.NewCoord(29, 18),
				"prim",
				"dijkstra",
			),
		},
	}

	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("#%d", i+1), func(t *testing.T) {
			t.Parallel()

			inputStream := bytes.NewBufferString(testCase.input)
			pres := presentation.New(inputStream, io.Discard)
			got, err := pres.ProcessInput()

			require.NoError(t, err, "input should get without error")
			require.Equal(t, testCase.expected, got, "values should be equal")
		})
	}
}

func TestProcessInputWithError(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		input string
	}{
		{
			input: "",
		},
		{
			input: "18\n18",
		},
		{
			input: "8\n8\n0\n0",
		},
		{
			input: "50\n50\n25\n0\n49",
		},
		{
			input: "40\n40\n20\n39\n39\n39\n",
		},
		{
			input: "5\n5\n0\n1\n4\n4\n2\n",
		},
		{
			input: "5 0",
		},
		{
			input: "5\n0",
		},
		{
			input: "5\n5\n0\n0\n3\n3",
		},
	}

	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("#%d", i+1), func(t *testing.T) {
			t.Parallel()

			inputStream := bytes.NewBufferString(testCase.input)
			pres := presentation.New(inputStream, io.Discard)
			_, err := pres.ProcessInput()

			require.Error(t, err, "input should return error")
		})
	}
}
