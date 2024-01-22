package constants

const INVALID_USER_REQUEST = "incorrect username or password"
const USERNAME_LENGTH_VALIDATION = "username and password must contain at least 6 characters"
const USERNAME_TAKEN = "username is already in use"
const INVALID_JWT = "invalid JWT in Access-Token header"
const INVALID_CLAIMS = "invalid claims in Access-Token header"
const UNAUTHORISED = "UNAUTHORISED"
const EMPTY_FIELD = "field must not be empty"
const DESCRIPTION_EMPTY = "description " + EMPTY_FIELD
const INVALID_DATE = "field must be of the format YYYY-MM-DD and be today or later"
const INVALID_TARGET_DATE = "targetDate " + INVALID_DATE
const TASK_ALREADY_DEFINED_ONE = "task with description: "
const TASK_ALREADY_DEFINED_TWO = " , and target date: "
const TASK_ALREADY_DEFINED_THREE = " is already defined"

func GetTaskAlreadyDefinedMessage(description, targetDate string) string {
	return TASK_ALREADY_DEFINED_ONE + description + TASK_ALREADY_DEFINED_TWO + targetDate + TASK_ALREADY_DEFINED_THREE
}
