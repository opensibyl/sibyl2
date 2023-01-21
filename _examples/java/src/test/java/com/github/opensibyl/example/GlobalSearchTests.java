package com.github.opensibyl.example;

import com.github.opensibyl.client.ApiClient;
import com.github.opensibyl.client.ApiException;
import com.github.opensibyl.client.Configuration;
import com.github.opensibyl.client.api.RegexQueryApi;
import com.github.opensibyl.client.model.Sibyl2FunctionWithPath;
import org.junit.Test;

import java.util.List;

public class GlobalSearchTests {
    @Test
    public void TestMain() {
        ApiClient defaultClient = Configuration.getDefaultApiClient();
        defaultClient.setBasePath(Constants.BASEURL);

        RegexQueryApi apiInstance = new RegexQueryApi(defaultClient);
        String repo = Constants.REPO; // String | repo
        String rev = Constants.REV; // String | rev
        try {
            List<Sibyl2FunctionWithPath> result = apiInstance.apiV1RegexFuncGet(
                    repo, rev, "name", ".*Handle.*");
            for (Sibyl2FunctionWithPath f : result) {
                System.out.println("function name with Handle: " + f.toString());
            }

        } catch (ApiException e) {
            System.err.println("Exception when calling DefaultApi#apiV1FileGet");
            System.err.println("Status code: " + e.getCode());
            System.err.println("Reason: " + e.getResponseBody());
            System.err.println("Response headers: " + e.getResponseHeaders());
            e.printStackTrace();
        }
    }
}
