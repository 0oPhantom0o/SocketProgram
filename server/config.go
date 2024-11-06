package main

//
//func config() {
//	data, err := os.ReadFile("config.yml")
//	if err != nil {
//		log.Fatalf("Error reading YAML file: %v", err)
//	}
//	fmt.Println("Raw YAML content:", string(data))
//	// Unmarshal the YAML data into a Config struct
//	err = yaml.Unmarshal(data, &conf)
//	if err != nil {
//		log.Fatalf("Error unmarshaling YAML: %v", err)
//	}
//	if conf.Port == "" || conf.Host == "" {
//		log.Println("One or more fields are not populated.")
//	}
//
//}
