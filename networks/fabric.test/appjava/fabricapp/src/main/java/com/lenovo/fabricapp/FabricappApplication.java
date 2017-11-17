package com.lenovo.fabricapp;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.web.bind.annotation.RestController;

@RestController
@SpringBootApplication
public class FabricappApplication {


	public static void main(String[] args) {


		SpringApplication.run(FabricappApplication.class, args);
	}
}
