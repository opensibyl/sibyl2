package com.github.opensibyl.example;

import com.github.opensibyl.client.ApiClient;
import com.github.opensibyl.client.ApiException;
import com.github.opensibyl.client.Configuration;
import com.github.opensibyl.client.StringUtil;
import com.github.opensibyl.client.api.BasicQueryApi;
import com.github.opensibyl.client.api.SignatureQueryApi;
import com.github.opensibyl.client.model.ObjectFunctionContextSlimWithSignature;
import com.github.opensibyl.client.model.ServiceFunctionContextChain;
import com.github.opensibyl.client.model.ServiceFunctionContextReverseChain;
import org.junit.Test;

import java.util.Arrays;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.stream.Collectors;

public class DiffTests {
    @Test
    public void testMain() {
        ApiClient defaultClient = Configuration.getDefaultApiClient();
        defaultClient.setBasePath(Constants.BASEURL);

        // assume that you have edited these files
        Map<String, List<Integer>> affectedMap = new HashMap<>();
        affectedMap.put("pkg/core/parser.go", Arrays.asList(4, 89, 90, 91, 92, 93, 94, 95, 96));
        affectedMap.put("pkg/core/unit.go", Arrays.asList(27, 28, 29));

        BasicQueryApi basicQueryApi = new BasicQueryApi(defaultClient);
        SignatureQueryApi signatureQueryApi = new SignatureQueryApi(defaultClient);
        affectedMap.forEach((k, v) -> {
            List<String> lines = v.stream().map(Object::toString)
                    .collect(Collectors.toList());
            try {
                List<ObjectFunctionContextSlimWithSignature> affectedFunctions = basicQueryApi
                        .apiV1FuncctxGet(
                                Constants.REPO,
                                Constants.REV,
                                k,
                                StringUtil.join(lines, ",")
                        );
                // to see their references
                affectedFunctions.forEach(each -> {
                    try {
                        ServiceFunctionContextChain reverseCallChain = signatureQueryApi.apiV1SignatureFuncctxRchainGet(
                                Constants.REPO,
                                Constants.REV,
                                each.getSignature(),
                                5);
                        // iterable tree-like object
                        reverseCallChain.getReverseCallChains().getChildren().forEach(eachChild -> {
                            eachChild.getChildren();
                            // ...
                        });
                    } catch (ApiException e) {
                        throw new RuntimeException(e);
                    }
                });
            } catch (ApiException e) {
                throw new RuntimeException(e);
            }
        });
    }
}
