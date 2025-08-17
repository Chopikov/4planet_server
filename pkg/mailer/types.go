package mailer

// EmailSubjects holds email subject configurations
type EmailSubjects struct {
	Verification  string
	PasswordReset string
}

// EmailURLs holds email URL configurations
type EmailURLs struct {
	BaseURL       string
	VerifyEmail   string
	ResetPassword string
}

// EmailTexts holds all email-related text configurations
type EmailTexts struct {
	Subjects EmailSubjects
	URLs     EmailURLs
	TeamName string
}
