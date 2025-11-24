package utils

func GetDisplayName(nickname, username string) string {
	if nickname != "" {
		return nickname
	}
	return username
}
