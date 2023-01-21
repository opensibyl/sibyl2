package com.github.opensibyl.example;

import com.github.opensibyl.client.ApiClient;
import com.github.opensibyl.client.ApiException;
import com.github.opensibyl.client.Configuration;
import com.github.opensibyl.client.api.DefaultApi;
import org.junit.Test;

import java.util.List;

public class SmokeTests {
    @Test
    public void testMain() {
        ApiClient defaultClient = Configuration.getDefaultApiClient();
        defaultClient.setBasePath(Constants.BASEURL);

        DefaultApi apiInstance = new DefaultApi(defaultClient);
        String repo = Constants.REPO; // String | repo
        String rev = Constants.REV; // String | rev
        try {
            List<String> result = apiInstance.apiV1FileGet(repo, rev);
            System.out.println(result);
        } catch (ApiException e) {
            System.err.println("Exception when calling DefaultApi#apiV1FileGet");
            System.err.println("Status code: " + e.getCode());
            System.err.println("Reason: " + e.getResponseBody());
            System.err.println("Response headers: " + e.getResponseHeaders());
            e.printStackTrace();
        }
    }
}
