package consts

const (
	Condition_EQUAL                 = "="
	Condition_NOT_EQUAL             = "!="
	Condition_GREATER_THAN          = ">"
	Condition_GREATER_THAN_OR_EQUAL = ">="
	Condition_LESS_THAN             = "<"
	Condition_LESS_THAN_OR_EQUAL    = "<="
	Condition_LIKE                  = "LIKE"
	Condition_NOT_LIKE              = "NOT LIKE"
	Condition_IN                    = "IN"
	Condition_NOT_IN                = "NOT IN"
	Condition_IS_NULL               = "IS NULL"
	Condition_IS_NOT_NULL           = "IS NOT NULL"
)

const (
	Join_INNER = "inner"
	Join_LEFT  = "left"
	Join_RIGHT = "right"
	Join_CROSS = "cross"
)

const (
	Join_Type_INNER = "INNER"
	Join_Type_LEFT  = "LEFT"
	Join_Type_RIGHT = "RIGHT"
	Join_Type_CROSS = "CROSS"
)

const (
	LogicalOperator_AND = iota
	LogicalOperator_OR
)

const (
	Order_ASC       = "ASC"
	Order_DESC      = "DESC"
	Order_FLAG_ASC  = true
	Order_FLAG_DESC = false
)

const (
	Lock_FOR_UPDATE = "FOR UPDATE"
	Lock_SHARE_MODE = "LOCK IN SHARE MODE"
)
