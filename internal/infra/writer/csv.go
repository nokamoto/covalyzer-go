package writer

import (
	"encoding/csv"
	"fmt"
	"os"
	"time"

	v1 "github.com/nokamoto/covalyzer-go/pkg/api/v1"
)

type CSVWriter struct{}

func NewCSVWriter() *CSVWriter {
	return &CSVWriter{}
}

func (c *CSVWriter) Write(config *v1.Config, data *v1.Covalyzer) error {
	var rows [][]string
	row := []string{"github", "repository"}
	for _, ts := range config.GetTimestamps() {
		d, err := time.Parse(time.RFC3339, ts)
		if err != nil {
			return fmt.Errorf("failed to parse timestamp: %s: %w", ts, err)
		}
		row = append(row, d.Format("2006-01-02"))
	}
	rows = append(rows, row)

	for _, repo := range data.GetRepositories() {
		row = []string{
			repo.GetRepository().GetGh(),
			fmt.Sprintf("%s/%s", repo.GetRepository().GetOwner(), repo.GetRepository().GetRepo()),
		}
		for _, coverage := range repo.GetCoverages() {
			row = append(row, fmt.Sprintf("%.1f", coverage.GetCover().GetTotal()))
		}
		rows = append(rows, row)
	}

	file, err := os.Create("covalyzer.csv")
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	w := csv.NewWriter(file)
	w.WriteAll(rows)
	return nil
}
