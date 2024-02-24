package writer

import (
	"encoding/csv"
	"fmt"
	"os"
	"time"

	v1 "github.com/nokamoto/covalyzer-go/pkg/api/v1"
)

type CSVWriter struct {
	file        string
	outlineFile string
}

func NewCSVWriter() *CSVWriter {
	return &CSVWriter{
		file:        "covalyzer.csv",
		outlineFile: "covalyzer-ginkgo-outline.csv",
	}
}

func (c *CSVWriter) header(config *v1.Config) ([]string, error) {
	row := []string{"github", "repository"}
	for _, ts := range config.GetTimestamps() {
		d, err := time.Parse(time.RFC3339, ts)
		if err != nil {
			return nil, fmt.Errorf("failed to parse timestamp: %s: %w", ts, err)
		}
		row = append(row, d.Format("2006-01-02"))
	}
	return row, nil
}

func (c *CSVWriter) writeFile(file string, rows [][]string) error {
	f, err := os.Create(file)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	w := csv.NewWriter(f)
	if err := w.WriteAll(rows); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}
	return nil
}

func (c *CSVWriter) repositoryColumns(repo *v1.RepositoryCoverages) []string {
	return []string{
		repo.GetRepository().GetGh(),
		fmt.Sprintf("%s/%s", repo.GetRepository().GetOwner(), repo.GetRepository().GetRepo()),
	}
}

func (c *CSVWriter) writeGo(config *v1.Config, data *v1.Covalyzer) error {
	var rows [][]string
	row, err := c.header(config)
	if err != nil {
		return err
	}
	rows = append(rows, row)

	for _, repo := range data.GetRepositories() {
		row = c.repositoryColumns(repo)
		for _, coverage := range repo.GetCoverages() {
			row = append(row, fmt.Sprintf("%.1f", coverage.GetCover().GetTotal()))
		}
		rows = append(rows, row)
	}

	return c.writeFile(c.file, rows)
}

func (c *CSVWriter) writeGinkgoOutline(config *v1.Config, data *v1.Covalyzer) error {
	var rows [][]string
	row, err := c.header(config)
	if err != nil {
		return err
	}
	rows = append(rows, row)

	for _, repo := range data.GetRepositories() {
		row = c.repositoryColumns(repo)
		for _, coverage := range repo.GetCoverages() {
			var nodes int32
			for _, ginkgo := range coverage.GetCover().GetGinkgo() {
				nodes += ginkgo.GetOutlineNodes()
			}
			row = append(row, fmt.Sprintf("%d", nodes))
		}
		rows = append(rows, row)
	}

	return c.writeFile(c.outlineFile, rows)
}

func (c *CSVWriter) Write(config *v1.Config, data *v1.Covalyzer) error {
	if err := c.writeGo(config, data); err != nil {
		return fmt.Errorf("failed to write go: %w", err)
	}
	if err := c.writeGinkgoOutline(config, data); err != nil {
		return fmt.Errorf("failed to write ginkgo outline: %w", err)
	}
	return nil
}
