package usercontext

var currentUserID uint = 1 // тестовый пользователь

func GetCurrentUserID() uint {
	return currentUserID
}

func GetCurrentUserIsModerator() bool {
	// Для тестов будем считать, что пользователь 1 является модератором (см. миграцию)
	return currentUserID == 1
}