package Error

var (
	InvalidInput_RequiredOption       = "option1, option2 or option3 are empty"
	InvalidInput_RequiredOptionWeight = "option1, option2 or option3 weight are empty"
	InvalidInput_NegativeOrZeroWeight = "weight must be positive"
	InvalidInput_InvalidSumOfWeight   = "the sum of the weights must be greater than 0 and less than 100"
	InvalidInput_DuplicateOption      = "options cannot be repeated"

	InvalidSignUp_UserNameCannotBeEmpty           = "username cannot be empty"
	InvalidSignUp_UserNameCannotContainSpaces     = "username cannot contain spaces"
	InvalidSignUp_UserNameAlreadyTaken            = "username is already taken"
	InvalidSignUp_PasswordTooShort                = "password is too short, the length must be above 8"
	InvalidSignUp_PasswordDiffFromConfirmPassword = "password is different from confirmation password"

	InvalidLogin_UserDoesNotExist        = "user does not exist"
	InvalidLogin_WrongUserNameOrPassword = "username or password is wrong"
	InvalidLogin_WrongPassword           = "password is wrong"

	InvalidApiKey         = "unable to read API key"
	InvalidSessionKey     = "unable to read session key"
	InvalidUser           = "user is unauthorized"
	HistoricalRecordEmpty = "historical records are empty"

	OutputError_CannotGetUserLocation = "unable to get user location"
	OutputError_CannotGetNearBy       = "unable to get nearBy"
	OutputError_AllZeroResults        = "the result is empty"
	OutputError_RequestError          = "system encounters request error"
)
