package mail

import "errors"

func findTemplate(template string) (string, error) {
	switch template {
	case "registration":
		return "Subject: Registration in Quote-service\nYou have been successfully registered", nil
	default:
		return "", errors.New("msg-template not found")
	}
}
