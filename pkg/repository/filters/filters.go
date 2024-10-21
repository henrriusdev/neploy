package filters

import (
	"time"

	"github.com/doug-martin/goqu/v9"
)

// SelectFilterBuilder defines a function signature for filters
type SelectFilterBuilder func(*goqu.SelectDataset) *goqu.SelectDataset

// ApplyFilters applies all filters to the query
func ApplyFilters(q *goqu.SelectDataset, filters ...SelectFilterBuilder) *goqu.SelectDataset {
	for _, filter := range filters {
		q = filter(q)
	}
	return q
}

// GenericColumnFilter handles filtering based on different types (string, bool, etc.)
func GenericColumnSelectFilter[T comparable](column string, value T, zeroValue T) SelectFilterBuilder {
	return func(q *goqu.SelectDataset) *goqu.SelectDataset {
		if value != zeroValue {
			return q.Where(goqu.I(column).Eq(value))
		}
		return q
	}
}

// DateRangeFilter creates a date range filter
// func DateRangeSelectFilter(startDate, endDate *model.Date, column string) SelectFilterBuilder {
// 	return func(q *goqu.SelectDataset) *goqu.SelectDataset {
// 		if startDate != nil && endDate != nil {
// 			dateRange := common.FormatDateRange(startDate.Time, endDate.Time)
// 			return q.Where(goqu.I(column).Gte(dateRange.StartDateStr), goqu.I(column).Lte(dateRange.EndDateStr))
// 		}
// 		return q
// 	}
// }

// TimeFilter filters based on start and/or end time
func TimeSelectFilter(startDate, endDate *time.Time, column string) SelectFilterBuilder {
	return func(q *goqu.SelectDataset) *goqu.SelectDataset {
		if startDate != nil {
			q = q.Where(goqu.I(column).Gte(*startDate))
		}
		if endDate != nil {
			q = q.Where(goqu.I(column).Lte(*endDate))
		}
		return q
	}
}

// LimitOffsetFilter applies limit and offset for pagination
func LimitOffsetFilter(limit, offset uint) SelectFilterBuilder {
	return func(q *goqu.SelectDataset) *goqu.SelectDataset {
		if limit > 0 {
			q = q.Limit(limit)
		}
		if offset > 0 {
			q = q.Offset(offset)
		}
		return q
	}
}

// IsFilter handles filtering for various types, including NULL values
func IsSelectFilter(column string, value interface{}) SelectFilterBuilder {
	return func(q *goqu.SelectDataset) *goqu.SelectDataset {
		switch v := value.(type) {
		case bool:
			return GenericColumnSelectFilter(column, v, false)(q)
		case nil:
			return q.Where(goqu.I(column).IsNull())
		case string:
			if v == "NOT NULL" {
				return q.Where(goqu.I(column).IsNotNull())
			} else {
				return GenericColumnSelectFilter(column, v, "")(q)
			}
		default:
			return q.Where(goqu.I(column).Eq(v))
		}
	}
}

// NotFilter negates the condition produced by another filter
func NotSelectFilter(filter SelectFilterBuilder) SelectFilterBuilder {
	return func(q *goqu.SelectDataset) *goqu.SelectDataset {
		// Apply the filter to get its expression
		clauses := filter(goqu.From()).GetClauses()

		// Check if the filter has a WHERE clause to negate
		if clauses.Where() != nil && clauses.Where().Expression() != nil {
			// Negate the WHERE clause
			return q.Where(goqu.L("NOT ?", clauses.Where().Expression()))
		}

		// Return the dataset unchanged if no WHERE clause exists
		return q
	}
}

// OrFilter takes the conditions produced by another filter and append them to an or
func OrSelectFilter(filters ...SelectFilterBuilder) SelectFilterBuilder {
	return func(q *goqu.SelectDataset) *goqu.SelectDataset {
		conditions := []goqu.Expression{}
		for _, filter := range filters {
			clauses := filter(goqu.From()).GetClauses()
			if clauses.Where() != nil && clauses.Where().Expression() != nil {
				conditions = append(conditions, clauses.Where().Expression())
			}
		}
		if len(conditions) > 0 {
			return q.Where(goqu.Or(conditions...))
		}
		return q
	}
}

// NumericComparisonFilter creates a filter for numeric comparisons
func NumericComparisonSelectFilter(column string, value int, op string) SelectFilterBuilder {
	return func(q *goqu.SelectDataset) *goqu.SelectDataset {
		switch op {
		case "lt":
			return q.Where(goqu.I(column).Lt(value))
		case "gt":
			return q.Where(goqu.I(column).Gt(value))
		case "eq":
			return q.Where(goqu.I(column).Eq(value))
		default:
			return q
		}
	}
}

// InFilter applies an "IN" condition on the specified column for a list of values
func InSelectFilter(column string, values []string) SelectFilterBuilder {
	return func(q *goqu.SelectDataset) *goqu.SelectDataset {
		if len(values) > 0 {
			return q.Where(goqu.I(column).In(values))
		}
		return q
	}
}

func CurrentDateSelectFilter(column string, op string) SelectFilterBuilder {
	return func(q *goqu.SelectDataset) *goqu.SelectDataset {
		switch op {
		case "lt":
			return q.Where(goqu.I(column).Lt(goqu.L("CURRENT_DATE")))
		case "gt":
			return q.Where(goqu.I(column).Gt(goqu.L("CURRENT_DATE")))
		default:
			return q
		}
	}
}

func ExtractYearSelectFilter(column string, year int) SelectFilterBuilder {
	return func(q *goqu.SelectDataset) *goqu.SelectDataset {
		return q.Where(goqu.L("EXTRACT(YEAR FROM ?) = ?", goqu.I(column), year))
	}
}

// UpdateFilterBuilder defines a function signature for filters
type UpdateFilterBuilder func(*goqu.UpdateDataset) *goqu.UpdateDataset

// ApplyUpdateFilters applies all filters to the update query
func ApplyUpdateFilters(q *goqu.UpdateDataset, filters ...UpdateFilterBuilder) *goqu.UpdateDataset {
	for _, filter := range filters {
		q = filter(q)
	}
	return q
}

// GenericColumnUpdateFilter handles filtering based on different types (string, bool, etc.)
func GenericColumnUpdateFilter[T comparable](column string, value T, zeroValue T) UpdateFilterBuilder {
	return func(q *goqu.UpdateDataset) *goqu.UpdateDataset {
		if value != zeroValue {
			return q.Where(goqu.I(column).Eq(value))
		}
		return q
	}
}

// DateRangeUpdateFilter creates a date range filter
// func DateRangeUpdateFilter(startDate, endDate *model.Date, column string) UpdateFilterBuilder {
// 	return func(q *goqu.UpdateDataset) *goqu.UpdateDataset {
// 		if startDate != nil && endDate != nil {
// 			dateRange := common.FormatDateRange(startDate.Time, endDate.Time)
// 			return q.Where(goqu.I(column).Gte(dateRange.StartDateStr), goqu.I(column).Lte(dateRange.EndDateStr))
// 		}
// 		return q
// 	}
// }

// TimeUpdateFilter filters based on start and/or end time
func TimeUpdateFilter(startDate, endDate *time.Time, column string) UpdateFilterBuilder {
	return func(q *goqu.UpdateDataset) *goqu.UpdateDataset {
		if startDate != nil {
			q = q.Where(goqu.I(column).Gte(*startDate))
		}
		if endDate != nil {
			q = q.Where(goqu.I(column).Lte(*endDate))
		}
		return q
	}
}

// IsUpdateFilter handles filtering for various types, including NULL values
func IsUpdateFilter(column string, value interface{}) UpdateFilterBuilder {
	return func(q *goqu.UpdateDataset) *goqu.UpdateDataset {
		switch v := value.(type) {
		case bool:
			return GenericColumnUpdateFilter(column, v, false)(q)
		case nil:
			return q.Where(goqu.I(column).IsNull())
		case string:
			if v == "NOT NULL" {
				return q.Where(goqu.I(column).IsNotNull())
			} else {
				return GenericColumnUpdateFilter(column, v, "")(q)
			}
		default:
			return q.Where(goqu.I(column).Eq(v))
		}
	}
}

// NotUpdateFilter negates the condition produced by another filter
func NotUpdateFilter(filter UpdateFilterBuilder) UpdateFilterBuilder {
	return func(q *goqu.UpdateDataset) *goqu.UpdateDataset {
		// Apply the filter to get its expression
		clauses := filter(q).GetClauses()

		// Check if the filter has a WHERE clause to negate
		if clauses.Where() != nil && clauses.Where().Expression() != nil {
			// Negate the WHERE clause
			return q.Where(goqu.L("NOT ?", clauses.Where().Expression()))
		}

		// Return the dataset unchanged if no WHERE clause exists
		return q
	}
}

// OrUpdateFilter takes the conditions produced by another filter and append them to an or
func OrUpdateFilter(filters ...UpdateFilterBuilder) UpdateFilterBuilder {
	return func(q *goqu.UpdateDataset) *goqu.UpdateDataset {
		conditions := []goqu.Expression{}
		for _, filter := range filters {
			clauses := filter(q).GetClauses()
			if clauses.Where() != nil && clauses.Where().Expression() != nil {
				conditions = append(conditions, clauses.Where().Expression())
			}
		}
		if len(conditions) > 0 {
			return q.Where(goqu.Or(conditions...))
		}
		return q
	}
}

// NumericComparisonUpdateFilter creates a filter for numeric comparisons
func NumericComparisonUpdateFilter(column string, value int, op string) UpdateFilterBuilder {
	return func(q *goqu.UpdateDataset) *goqu.UpdateDataset {
		switch op {
		case "lt":
			return q.Where(goqu.I(column).Lt(value))
		case "gt":
			return q.Where(goqu.I(column).Gt(value))
		default:
			return q
		}
	}
}

// InUpdateFilter applies an "IN" condition on the specified column for a list of values
func InUpdateFilter(column string, values []string) UpdateFilterBuilder {
	return func(q *goqu.UpdateDataset) *goqu.UpdateDataset {
		if len(values) > 0 {
			return q.Where(goqu.I(column).In(values))
		}
		return q
	}
}

func CurrentDateUpdateFilter(column string, op string) UpdateFilterBuilder {
	return func(q *goqu.UpdateDataset) *goqu.UpdateDataset {
		switch op {
		case "lt":
			return q.Where(goqu.I(column).Lt(goqu.L("CURRENT_DATE")))
		case "gt":
			return q.Where(goqu.I(column).Gt(goqu.L("CURRENT_DATE")))
		default:
			return q
		}
	}
}
