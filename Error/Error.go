package Error

var (
	InvalidInput_RequiredOption       = "option1, option2 or option3 are empty"
	InvalidInput_RequiredOptionWeight = "option1, option2 or option3 weight are empty"
	InvalidInput_NegativeOrZeroWeight = "weight must be positive"
	InvalidInput_InvalidSumOfWeight   = "the sum of the weights must be greater than 0 and less than 100"
	InvalidInput_DuplicateOption      = "options cannot be repeated"

	InvalidLogin_UserDoesNotExist        = "user does not exist"
	InvalidLogin_WrongUserNameOrPassword = "username or password is wrong"
	InvalidLogin_WrongPassword           = "password is wrong"

	InvalidApiKey = "unable to read API key"

	CannotGetUserLocation = "unable to get user location"
)
