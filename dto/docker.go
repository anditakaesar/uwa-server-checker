package dto

import "strings"

type Container struct {
	ID     string
	Names  []string
	State  string
	Status string
}

func (c *Container) GetNames() string {
	strippedNames := make([]string, 0, len(c.Names))
	for _, name := range c.Names {
		strippedNames = append(strippedNames, name[1:])
	}

	return strings.Join(strippedNames, ", ")
}
