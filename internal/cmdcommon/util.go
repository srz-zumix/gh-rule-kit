package cmdcommon

import (
	"fmt"
	"strconv"
)

// parseRulesetID parses a single ruleset ID positional argument.
func parseRulesetID(arg string) (int64, error) {
	id, err := strconv.ParseInt(arg, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid ruleset ID: %w", err)
	}
	return id, nil
}

// ParseRulesetIDs parses positional ruleset ID arguments into int64s.
func ParseRulesetIDs(args []string) ([]int64, error) {
	ids := make([]int64, 0, len(args))
	for _, idStr := range args {
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid ruleset ID '%s': %w", idStr, err)
		}
		ids = append(ids, id)
	}
	return ids, nil
}
