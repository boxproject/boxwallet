package errors

import "errors"

var (
	ERR_NIL_REFERENCE = errors.New("Null reference error.")

	ERR_COIN_PREC_OVERFLOW = errors.New("The number exceeds the minimum accuracy.")

	ERR_DIFF_UNIT = errors.New("Different units cannot be converted.")

	ERR_TX_WITHOUT_SGIN = errors.New("Tx is not signed.")

	ERR_PARAM_NOT_VALID = errors.New("Parameter not valid.")

	ERR_NOT_ENOUGH_COIN = errors.New("Not enough coin.")

	ERR_DATA_EXISTS = errors.New("Data already exists.")

	ERR_TX_PENDING = errors.New("pending ....")

	ERR_ADDRESS_QUEUE_BLOCKED = errors.New("Current address is blocked,cannot create a second transaction.")

	ERR_TX_END_NOT_NORMAL = errors.New("Transaction abort with exception.")

	ERR_KEY_OVERFLOW = errors.New("Can't overflow.")

	ERR_TIME_OUT = errors.New("time out.")

	ERR_TX_OUT_INDEX_OVERFLEW = errors.New("tx out index out of range")

	ERR_PIPELINE_DATA_ILLEGAL = errors.New("Internal identification of data illegal")
)
