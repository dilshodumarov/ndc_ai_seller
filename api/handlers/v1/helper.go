package v1


import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"sugurta/internal/entity"
)

const pythonBaseURL = "http://telegram-fastapi:8000"
//const pythonBaseURL = "http://localhost:8000"
// SendTelegramCode sends a phone number to the Python backend and returns the response
func SendTelegramCode(phone entity.PhoneNumber) (*entity.BotIntegrationResponse, error) {
	body, err := json.Marshal(phone)
	if err != nil {
		log.Printf("❌ JSON marshal xatolik: %v", err)
		return nil, err
	}

	resp, err := http.Post(pythonBaseURL+"/login/send-code", "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Printf("❌ HTTP POST xatolik (send-code): %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	var Resp entity.BotIntegrationResponse
	if err := json.NewDecoder(resp.Body).Decode(&Resp); err != nil {
		log.Printf("❌ JSON decode xatolik: %v", err)
		return nil, err
	}

	if Resp.Code != 0 {
		log.Printf("⚠️ Python botdan kelgan xatolik kodi: %d, message: %s", Resp.Code, Resp.Message)
		return &Resp, nil
	}

	log.Printf("✅ Kod yuborildi: %+v", Resp)
	return &Resp, nil
}


// SendTelegramVerify verifies the code and optional password with the Python backend
func SendTelegramVerify(input entity.CodeInput) (*entity.BotIntegrationResponse, error) {
	body, err := json.Marshal(input)
	if err != nil {
		log.Printf("❌ JSON marshal xatolik (verify): %v", err)
		return nil, err
	}

	resp, err := http.Post(pythonBaseURL+"/login/verify", "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Printf("❌ HTTP POST xatolik (verify): %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	var Resp entity.BotIntegrationResponse
	if err := json.NewDecoder(resp.Body).Decode(&Resp); err != nil {
		log.Printf("❌ JSON decode xatolik (verify): %v", err)
		return nil, err
	}

	if Resp.Code != 0 {
		log.Printf("⚠️ Python botdan verify xatolik: %d, message: %s", Resp.Code, Resp.Message)
		return &Resp, nil
	}

	log.Printf("✅ Kod tasdiqlandi: %+v", Resp)
	return &Resp, nil
}


// SendTelegramMessage sends a message request to the Python backend
func SendTelegramMessage(msg entity.MessageRequest) (*entity.BotIntegrationResponse, error) {
	body, err := json.Marshal(msg)
	if err != nil {
		log.Printf("❌ JSON marshal xatolik (send-message): %v", err)
		return nil, err
	}

	resp, err := http.Post(pythonBaseURL+"/send-message/", "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Printf("❌ HTTP POST xatolik (send-message): %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	var Resp entity.BotIntegrationResponse
	if err := json.NewDecoder(resp.Body).Decode(&Resp); err != nil {
		log.Printf("❌ JSON decode xatolik (send-message): %v", err)
		return nil, err
	}
	if Resp.Code != 0 {
		log.Printf("⚠️ Python botdan message xatolik: %d, message: %s", Resp.Code, Resp.Message)
		return &Resp, nil
	}

	log.Printf("✅ Xabar yuborildi: %+v", Resp)
	return &Resp, nil
}


func SendTelegramStartSession(input entity.PhoneNumber) (*entity.BotIntegrationResponse, error) {
	body, err := json.Marshal(input)
	if err != nil {
		log.Printf("❌ JSON marshal xatolik (start): %v", err)
		return nil, err
	}

	resp, err := http.Post(pythonBaseURL+"/session/start", "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Printf("❌ HTTP POST xatolik (start): %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	var Resp entity.BotIntegrationResponse
	if err := json.NewDecoder(resp.Body).Decode(&Resp); err != nil {
		log.Printf("❌ JSON decode xatolik (start): %v", err)
		return nil, err
	}

	if Resp.Code != 0 {
		log.Printf("⚠️ Python botdan start xatolik: %d, message: %s", Resp.Code, Resp.Message)
		return &Resp, nil
	}

	log.Printf("✅ Telegram sessiya boshlandi: %+v", Resp)
	return &Resp, nil
}


func SendTelegramStopSession(input entity.PhoneNumber) (*entity.BotIntegrationResponse, error) {
	body, err := json.Marshal(input)
	if err != nil {
		log.Printf("❌ JSON marshal xatolik (stop): %v", err)
		return nil, err
	}

	resp, err := http.Post(pythonBaseURL+"/session/stop", "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Printf("❌ HTTP POST xatolik (stop): %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	var Resp entity.BotIntegrationResponse
	if err := json.NewDecoder(resp.Body).Decode(&Resp); err != nil {
		log.Printf("❌ JSON decode xatolik (stop): %v", err)
		return nil, err
	}

	if Resp.Code != 0 {
		log.Printf("⚠️ Python botdan stop xatolik: %d, message: %s", Resp.Code, Resp.Message)
		return &Resp, nil
	}

	log.Printf("✅ Telegram sessiya to‘xtatildi: %+v", Resp)
	return &Resp, nil
}
