package dao

// GetSetting 返回指定 key 的设置值，如果不存在返回空字符串 + false
//func GetSetting(key string) (string, bool, error) {
//	var setting models.UserSetting
//	err := repo.DB.Where("key = ?", key).First(&setting).Error
//	if errors.Is(err, gorm.ErrRecordNotFound) {
//		return "", false, nil
//	}
//	if err != nil {
//		return "", false, err
//	}
//	return setting.Value, true, nil
//}

// SetSetting 设置或更新某项配置
//func SetSetting(key string, value string) error {
//	var setting models.UserSetting
//	err := repo.DB.Where("key = ?", key).First(&setting).Error
//	if errors.Is(err, gorm.ErrRecordNotFound) {
//		// 新建
//		return repo.DB.Create(&models.UserSetting{
//			Key:   key,
//			Value: value,
//		}).Error
//	} else if err != nil {
//		return err
//	}
//	// 更新已有
//	setting.Value = value
//	return repo.DB.Save(&setting).Error
//}
//
//// GetSettingWithDefault 如果未设置，则返回默认值
//func GetSettingWithDefault(key string, defaultVal string) (string, error) {
//	val, exists, err := GetSetting(key)
//	if err != nil {
//		return "", err
//	}
//	if !exists {
//		return defaultVal, nil
//	}
//	return val, nil
//}
