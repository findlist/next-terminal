package api

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"io"
	"next-terminal/server/common/maps"
	"next-terminal/server/service"
)

type LicenseApi struct {
}

func (api LicenseApi) MachineIdLicense(c echo.Context) error {
	machineCode := ""
	machineCode = service.LicenseService.GetMachineId()
	return Success(c, machineCode)
}
func (api LicenseApi) License(c echo.Context) error {
	license, err := service.LicenseService.GetLicense()
	if err != nil {
		return err
	}

	return Success(c, license)
}
func (api LicenseApi) LicenseMi(c echo.Context) error {
	m := maps.Map{}
	if err := c.Bind(&m); err != nil {
		return err
	}
	// 将 JSON 数据编码为字节切片
	jsonData, err := json.Marshal(m)
	if err != nil {
		return err
	}
	if jsonData == nil {
		return nil
	}
	fmt.Println("加密前:", string(jsonData))
	// 使用 AES 加密 JSON 数据
	encryptedData, err := encryptAES([]byte(jsonData), []byte("abcdefghijklmnopqrstuvwxyzwwwwww"))
	if err != nil {
		fmt.Println("Error encrypting data:", err)
		return nil
	}

	// 打印加密后的数据
	fmt.Println("Encrypted data:", encryptedData)

	// 将密文转换为Base64字符串以便存储或传输
	cipherTextBase64 := base64.StdEncoding.EncodeToString(encryptedData)
	fmt.Println("Encrypted data (Base64):", cipherTextBase64)

	return Success(c, cipherTextBase64)
}

// 使用 AES 加密数据
func encryptAES(data []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	ciphertext := make([]byte, aes.BlockSize+len(data))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], data)

	return ciphertext, nil
}
func (api LicenseApi) LicenseCreate(c echo.Context) error {
	// 解析请求中的 JSON 数据到 map 中
	m := maps.Map{}
	if err := c.Bind(&m); err != nil {
		return err
	}
	// 获取 id 参数的值
	licenseValue, ok := m["license"]
	if !ok {
		// 如果不存在 LicenseValue 参数，则返回错误或者采取其他处理方式
		return errors.New("License parameter not found")
	}
	// 将 license 参数的值转换为字符串
	licenseStr, ok := licenseValue.(string)
	if !ok {
		// 如果无法将值转换为字符串，则返回错误或者采取其他处理方式
		return errors.New("license parameter is not a string")
	}
	// 调用 LicenseService 中的 Create 方法来创建许可证
	license, err := service.LicenseService.Create(licenseStr)
	if err != nil {
		return err
	}

	// 返回成功的响应
	return Success(c, license)
}
func (api LicenseApi) LicenseUpdate(c echo.Context) error {
	id := c.Param("id")
	m := maps.Map{}
	if err := c.Bind(&m); err != nil {
		return err
	}
	if err := service.LicenseService.UpdateById(id, m); err != nil {
		return err
	}
	return Success(c, nil)
}
func (api LicenseApi) LicenseDelete(c echo.Context) error {
	id := c.Param("id")
	if err := service.LicenseService.DeleteById(id); err != nil {
		return err
	}

	return Success(c, nil)
}
