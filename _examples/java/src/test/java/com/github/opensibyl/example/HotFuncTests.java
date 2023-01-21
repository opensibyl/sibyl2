package com.github.opensibyl.example;

import com.github.opensibyl.client.ApiClient;
import com.github.opensibyl.client.ApiException;
import com.github.opensibyl.client.Configuration;
import com.github.opensibyl.client.api.ReferenceQueryApi;
import com.github.opensibyl.client.api.RegexQueryApi;
import com.github.opensibyl.client.model.Sibyl2FunctionContextSlim;
import com.github.opensibyl.client.model.Sibyl2FunctionWithPath;
import org.junit.Test;

import java.util.List;

public class HotFuncTests {
    @Test
    public void TestMain() {
        ApiClient defaultClient = Configuration.getDefaultApiClient();
        defaultClient.setBasePath(Constants.BASEURL);

        ReferenceQueryApi apiInstance = new ReferenceQueryApi(defaultClient);
        String repo = Constants.REPO; // String | repo
        String rev = Constants.REV; // String | rev
        try {
            List<Sibyl2FunctionContextSlim> funcctxs = apiInstance.apiV1ReferenceCountFuncctxGet(repo, rev, 10, 100);
            for (Sibyl2FunctionContextSlim each : funcctxs) {
                System.out.printf("f: %s called: %s%n", each.getName(), each.getReverseCalls());
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
