package com.lenovo.fabricapp.controller;


import com.lenovo.fabricapp.respository.InvokeChainCode;
import net.sf.json.JSONObject;
import org.springframework.stereotype.Controller;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestMethod;
import org.springframework.web.bind.annotation.ResponseBody;

import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;
import java.io.IOException;

@Controller
@RequestMapping(value="/lenovo/")
public class LenovoController {
    @RequestMapping(value = "get", method = RequestMethod.GET)
    public @ResponseBody
    String getInfo(String key, HttpServletRequest request, HttpServletResponse response) throws IOException {

        if (key == null ) {
            return "error: key is null";
        }

        String[] args = new String[]{"getInsurance", key};

        String result = null;
        try {
            InvokeChainCode invoke = new InvokeChainCode(args);
            result = invoke.invoke();
        } catch (Exception e) {
            e.printStackTrace();
        }
        System.out.println("result:" + result);
        return result;
    }

    @RequestMapping(value = "put", method = RequestMethod.POST)
    public @ResponseBody String postInfo(@RequestBody String body, HttpServletRequest request, HttpServletResponse response) throws IOException {
        System.out.println("body:" + body);
        JSONObject requestcontent = JSONObject.fromObject(body);
        System.out.println("requestcontent:" + requestcontent);
//        String ID = requestcontent.getString("ID");
//        String ProductName = requestcontent.getString("ProductName");
//        String ProductType = requestcontent.getString("ProductType”");
//        String OrganizationID = requestcontent.getString("OrganizationID");
//        String Portion = requestcontent.getString("Portion");

//        return "ok";

        String data = "{\n" +
                "\t\"PolicyNo\":\"123\",\n" +
                "\t\"InsurantID\":\"456\",\n" +
                "\t\"ServiceAgreementHASH\":\"ABC\",\n" +
                "\t\"IMEINo\":\"mobile1\",\n" +
                "\t\"ActivateStoreID\":\"no1\",\n" +
                "\t\"VerifyResult\":\"pass\",\n" +
                "\t\"SignDate\":\"today\",\n" +
                "\t\"EffectiveDate\":\"tomorrow\",\n" +
                "\t\"ExpirationDate\":\"the day after tomorrow\",\n" +
                "\t\"ClaimFlag\":\"true\",\n" +
                "\t\"ClaimDate\":\"now\",\n" +
                "\t\"ClaimAmount\":\"1000\",\n" +
                "\t\"ModifyFlag\":\"yes\",\n" +
                "\t\"ExtraData\":[\n" +
                "\t\t{\"Data\":\"hello\"},\n" +
                "\t\t{\"Data\":\"world\"}\n" +
                "\t\t]\n" +
                "}";

        String[] args = new String[]{"postInsurance", body};

        String result = null;
        try {
            InvokeChainCode invoke = new InvokeChainCode(args);
            result = invoke.invoke();
        } catch (Exception e) {
            e.printStackTrace();
        }
        return result;

    }


}
