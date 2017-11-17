package com.lenovo.fabricapp.controller;

import com.lenovo.fabricapp.respository.*;
import org.springframework.stereotype.Controller;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestMethod;
import org.springframework.web.bind.annotation.ResponseBody;
import org.springframework.web.bind.annotation.RestController;

/**
 * Created by dzkan on 2016/3/8.
 */
@RestController
public class MainController {

    // 自动装配数据库接口，不需要再写原始的Connection来操作数据库
//    @Autowired
//    UserRepository userRepository;

    @RequestMapping(value = "/", method = RequestMethod.GET)
    public String index() {
        return "Transactionabc";
    }
    @RequestMapping(value = "/install", method = RequestMethod.GET)
    public @ResponseBody
    String Inatall() {

        try {
            SetupUsers setupUsers=new SetupUsers();
            setupUsers.setup();
            ConstructChannel constructChannel=new ConstructChannel();
            constructChannel.constructchannel();
//            DeployChaincode deploy = new DeployChaincode();
//            deploy.install();
        }catch (Exception e){
            e.printStackTrace();
        }
        return "success";
    }
    @RequestMapping(value = "/reinstall", method = RequestMethod.GET)
    public @ResponseBody
    String ReInatall() {

        try {
            getUsers getusers=new getUsers();
            getusers.getusers();
            ReconstructChannel constructChannel=new ReconstructChannel();
            constructChannel.reconstructchannel();
            DeployChaincode deploy = new DeployChaincode();
            deploy.install();
        }catch (Exception e){
            e.printStackTrace();
        }
        return "success";
    }


}
