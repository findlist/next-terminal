package service

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"next-terminal/server/common/maps"
	"next-terminal/server/env"
	"next-terminal/server/model"
	"next-terminal/server/repository"
	"next-terminal/server/utils"
	"os/exec"
	"runtime"
	"strings"
)

var LicenseService = new(licenseService)

type licenseService struct {
	baseService
}

func (s licenseService) GetLicense() (model.License, error) {
	license, err := repository.LicenseRepository.FindLicense(context.TODO())
	if err != nil {
		return model.License{}, err
	}
	total, err := repository.AssetRepository.Count(context.TODO())
	license.SurplusAssets = int64(license.Asset) - total
	// 返回结果
	return license, nil
}
func (s licenseService) Create(licenseData string) (model.License, error) {
	// 密文数据经过Base64编码，需要解码为原始的密文数据
	decodedCipherText, err := base64.StdEncoding.DecodeString(licenseData)
	if err != nil {
		// 处理解码错误
		return model.License{}, err
	}
	// 使用 AES 解密数据
	decryptedData, err := decryptAES([]byte(decodedCipherText), []byte("abcdefghijklmnopqrstuvwxyzwwwwww"))
	if err != nil {
		fmt.Println("Error decrypting data:", err)
		return model.License{}, nil
	}

	// 打印解密后的数据
	fmt.Println("Decrypted data:", string(decryptedData))
	var license model.License
	if err := json.Unmarshal(decryptedData, &license); err != nil {
		fmt.Println("Error parsing license data:", err)
		return license, err
	}
	fmt.Println("license:", license)
	//比较机器码
	machineCode := ""
	machineCode = LicenseService.GetMachineId()
	if license.MachineID != machineCode {
		return model.License{}, errors.New("机器码异常")
	}
	//判断是否存在
	getLicense, err := LicenseService.GetLicense()
	if err != nil {
		fmt.Println("Error:", err)

		// 设置许可证的 ID
		license.ID = utils.UUID()

		// 开始事务
		err = s.Transaction(context.Background(), func(ctx context.Context) error {
			// 将许可证保存到数据库
			if err := repository.LicenseRepository.Create(ctx, &license); err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			fmt.Println("Error creating license:", err)
			return model.License{}, err
		}

		return model.License{}, err
	}
	license.ID = getLicense.ID
	err = s.Transaction(context.Background(), func(ctx context.Context) error {
		if err := repository.LicenseRepository.UpdateById(ctx, &license, getLicense.ID); err != nil {
			return err
		}
		return nil
	})
	return license, nil
}
func (s licenseService) UpdateById(id string, m maps.Map) error {
	data, err := json.Marshal(m)
	if err != nil {
		return err
	}
	var item model.License
	if err := json.Unmarshal(data, &item); err != nil {
		return err
	}

	return s.Transaction(context.Background(), func(ctx context.Context) error {
		if err := repository.LicenseRepository.UpdateById(ctx, &item, id); err != nil {
			return err
		}
		return nil
	})
}
func (s licenseService) DeleteById(id string) error {
	return env.GetDB().Transaction(func(tx *gorm.DB) error {
		c := s.Context(tx)
		// 删除资产
		if err := repository.LicenseRepository.DeleteById(c, id); err != nil {
			return err
		}
		return nil
	})
}

func (s licenseService) GetMachineId() (c string) {
	macAddress, err := getMACAddress()
	if err != nil {
		panic(err)
	}
	fmt.Println("MAC地址:", macAddress)

	diskSerialNumber, err := getDiskSerialNumber()
	if err != nil {
		panic(err)
	}
	// 获取硬盘序列号
	fmt.Println("硬盘序列号:", diskSerialNumber)

	osInfo := getOSInfo()
	// 获取操作系统信息
	fmt.Println("操作系统:", osInfo)

	// 使用获取的信息生成机器码
	machineInfo := fmt.Sprintf("%s|%s|%s", macAddress, diskSerialNumber, osInfo)

	// 计算机器码的SHA-1哈希值
	hasher := sha1.New()
	hasher.Write([]byte(machineInfo))
	machineCode := fmt.Sprintf("%x", hasher.Sum(nil))

	fmt.Println("Machine Code:", machineCode)

	return machineCode
}
func getMACAddress() (string, error) {
	var macAddr string
	if runtime.GOOS == "windows" {
		cmd := exec.Command("getmac")
		output, err := cmd.Output()
		if err != nil {
			return "", err
		}
		macAddr = strings.Split(string(output), "\n")[3]
	} else {
		cmd := exec.Command("/sbin/ifconfig", "-a")
		output, err := cmd.Output()
		if err != nil {
			return "", err
		}
		macAddr = strings.Split(string(output), "\n")[0]
	}
	return macAddr, nil
}

func getDiskSerialNumber() (string, error) {
	cmd := exec.Command("wmic", "diskdrive", "get", "serialnumber")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	serialNumber := strings.Split(string(output), "\n")[1]
	return serialNumber, nil
}

func getOSInfo() string {
	return runtime.GOOS
}

// generateAESKey 生成AES密钥
func generateAESKey(keySize int) ([]byte, error) {
	key := make([]byte, keySize)
	_, err := rand.Read(key)
	if err != nil {
		return nil, err
	}
	return key, nil
}

// 使用 AES 解密数据
func decryptAES(ciphertext []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	if len(ciphertext) < aes.BlockSize {
		return nil, fmt.Errorf("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	return ciphertext, nil
}
