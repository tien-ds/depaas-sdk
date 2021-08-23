import com.cos.ETH;
import com.cos.Oracle;
import com.cos.Signer;
import com.ds.DSErc20;

public class Test {
    static {
        System.loadLibrary("share");
    }

    public static void main(String[] args) {
        Signer.initClient();
        // Oracle.gatewayInit();
        // testEth();

        // testloom();

        // testds();

        // testCoin();

    }

    private static void testds() {

        DSErc20 ds= DSErc20.Of("0xc023e763954Ed2919FBC7F5fdE4502dA4A173cf8");

//         System.out.println(ds.stackOf("/ic1NwEWmdWx/BPDTI6dZSjlS88="));
//         System.out.println(Signer.getBalance("0xEAA9dF8450a74F2A75713Ed4475dA5cEDE4029A4",
//                                              "/ic1NwEWmdWx/BPDTI6dZSjlS88="));
//
//         System.out.println(ds.penalty("qV7nCltu0cuAzptROJbr8xQGWk6uuJBXHDkS3ylkRlaRmtjq7e7NSMQlguMj2XzBw3SsVn3GXcNIdCCRlaGZ6w==",
//         "/ic1NwEWmdWx/BPDTI6dZSjlS88=",
//         "100000000"));
//
//         System.out.println(Signer.getBalance("0xEAA9dF8450a74F2A75713Ed4475dA5cEDE4029A4",
//                                              "/ic1NwEWmdWx/BPDTI6dZSjlS88="));
//
        System.out.println(ds.stackOf("/ic1NwEWmdWx/BPDTI6dZSjlS88="));
//
        System.out.println(ds.owner());
    }

    private static void testloom() {
        // loom
        System.out.println(Signer.genKey());
        System.out.println(
                Signer.sign("HFEv1N8WBMtWS3PzM46PbGxQuvHaobbb/8Fcgi6UTFP1APdieOAIm2bMqa4MaEhFbPWO/Hi8lxzqZPCf9dY0Yg==",
                        "szb".getBytes()));
        System.out.println(Signer.getMapper("{\"type\":0,\"hexAddr\":\"0xdb834d1f5baf312424fe3003524e2f5a52bf15b2\"}"));
        System.out.println(Signer.tokenInfo("0xed36585a83f41eb3b37b061bd1a297605a890513"));
        System.out.println(Signer.getBalance("0xed36585a83f41eb3b37b061bd1a297605a890513",
                "0xdb834d1f5baf312424fe3003524e2f5a52bf15b2"));
    }

    private static void testEth() {
        // eth
        System.out.println(ETH.getTokenBalance("http://122.5.20.58:25545", "0x0A37B2fdCC10C80Da61804EB7Da01b5AD3C66ecC",
                "0x21d1acf804e0423431452b161d39cb8f8fea6b0d"));
        new Oracle("0x0a37b2fdcc10c80da61804eb7da01b5ad3c66ecc", "0xc1b174f1bc70172c911df8a01d0e0d98129c7517");

        System.out.println(Signer.hexToAddr("0x3ef8b43b5c1577ee37cf6c314642f521e6451824"));
        
    }

    private static void testCoin() {
        System.out.println();
    }

}