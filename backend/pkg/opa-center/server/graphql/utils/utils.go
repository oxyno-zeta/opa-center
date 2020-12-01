package utils

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/oxyno-zeta/opa-center/pkg/opa-center/common/errors"
	"github.com/oxyno-zeta/opa-center/pkg/opa-center/database/pagination"
)

const maxPageSize = 50
const paginationIDPrefix = "paginate"
const relayIDSplitSize = 2

func ToIDRelay(prefix, id string) string {
	return base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", prefix, id)))
}

func FromIDRelay(relayID, prefix string) (string, error) {
	// Base64 decode
	idBb, err := base64.StdEncoding.DecodeString(relayID)
	// Check error
	if err != nil {
		return "", errors.NewInvalidInputErrorWithError(err)
	}

	idContent := string(idBb)
	// Split
	sp := strings.Split(idContent, ":")
	if len(sp) != relayIDSplitSize {
		return "", errors.NewInvalidInputError("format error on relay token")
	}
	// Check that first item of split is a good
	if sp[0] != prefix {
		return "", errors.NewInvalidInputError("invalid relay prefix")
	}

	return sp[1], nil
}

func FormatTime(ti time.Time) string {
	return ti.Format(time.RFC3339)
}

func GetPaginateCursor(tableIndex, skip int) string {
	return ToIDRelay(paginationIDPrefix, fmt.Sprintf("%d", tableIndex+skip+1))
}

func GetPageInfo(startCursor, endCursor string, p *pagination.PageOutput) *PageInfo {
	var res PageInfo

	// Check if start cursor exists
	if startCursor != "" {
		res.StartCursor = &startCursor
	}

	// Check if end cursor exists
	if endCursor != "" {
		res.EndCursor = &endCursor
	}

	// Check if paginator exists
	if p != nil {
		res.HasNextPage = p.HasNext
		res.HasPreviousPage = p.HasPrevious
	}

	return &res
}

func GetPageInput(after *string, before *string, first *int, last *int) (*pagination.PageInput, error) {
	// Check if all cursors are present together
	if after != nil && before != nil {
		return nil, errors.NewInvalidInputError("after and before can't be present together at the same time")
	}
	// Check if first and last are present together
	if first != nil && last != nil {
		return nil, errors.NewInvalidInputError("first and last can't be present together at the same time")
	}
	// Check before and last
	if before != nil && last == nil {
		return nil, errors.NewInvalidInputError("before must be used with last element")
	}
	// Check before and last case 2
	if (before == nil || *before == "") && last != nil {
		return nil, errors.NewInvalidInputError("last must be used with before element")
	}
	// Check first and after
	if after != nil && first == nil {
		return nil, errors.NewInvalidInputError("first must be used with after element")
	}
	// Check if last is positive
	if last != nil && *last <= 0 {
		return nil, errors.NewInvalidInputError("last must be > 0")
	}
	// Check if first is positive
	if first != nil && *first <= 0 {
		return nil, errors.NewInvalidInputError("first must be > 0")
	}

	// Create parginator input
	var res pagination.PageInput

	// Before case
	if before != nil && *before != "" {
		i, err := parsePaginateCursor(*before)
		// Check error
		if err != nil {
			return nil, err
		}

		res.Skip = i - *last - 1
		res.Limit = *last

		// Check if skip is positive
		if res.Skip < 0 {
			res.Skip = 0
		}
	}

	// After case
	if after != nil && *after != "" {
		i, err := parsePaginateCursor(*after)
		// Check error
		if err != nil {
			return nil, err
		}

		res.Skip = i
		res.Limit = *first
	}

	// First not null and after is
	if (after == nil || *after == "") && first != nil {
		res.Limit = *first
	}

	// Check limit
	if res.Limit > maxPageSize {
		return nil, errors.NewInvalidInputError(fmt.Sprintf("first or last is too big, maximum is %d", maxPageSize))
	}

	// Set default default limit
	if res.Limit == 0 {
		res.Limit = 10
	}

	return &res, nil
}

func parsePaginateCursor(cursorB64 string) (int, error) {
	val, err := FromIDRelay(cursorB64, paginationIDPrefix)
	// Check error
	if err != nil {
		return 0, err
	}
	// Parse cursor int
	res, err := strconv.Atoi(val)
	// Check error
	if err != nil {
		return 0, err
	}
	// Check if cursor is positive
	if res < 0 {
		return 0, errors.NewInvalidInputError("cursor pagination must be positive")
	}

	return res, nil
}
