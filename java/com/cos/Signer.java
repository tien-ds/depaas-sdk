package com.cos;

public class Signer {

    public static native String genKey();

    public static native String genKeyFromSeed(String seed);

    public static native String getBalance(String contractAddr, String hexAddr);

    /**
     * { "contractAddr": "0xdb834d1f5baf312424fe3003524e2f5a52bf15b2",
     * "to":"0xdb834d1f5baf312424fe3003524e2f5a52bf15b2", "base64PrivateKey": "",
     * "amount": 10000 }
     **/
    public static native String transfer(String sJson);

    /**
    * support amount string
    *
    * { "contractAddr": "0xdb834d1f5baf312424fe3003524e2f5a52bf15b2",
    * "to":"0xdb834d1f5baf312424fe3003524e2f5a52bf15b2", "base64PrivateKey": "",
    * "amount": "10000" }
    **/
    public static native String transfer1(String sJson);

    /**
     * {"type":0,"hexAddr":"0xdb834d1f5baf312424fe3003524e2f5a52bf15b2"} type 0 eth,
     * 1 default
     **/
    public static native String getMapper(String sJson);

    // account mapping
    public static native String mapAccount(String mainNetPrivateKey, String base64PrivateKey);

    // get token info
    public static native String tokenInfo(String contractAddr);

    public static native String hexToAddr(String hex);

    public static native String addrToHex(String hex);

    public static native String sign(String base64PrivateKey, byte[] data);

    public static native String privateKeyToAddr(String base64PrivateKey);

    //caller string, contract string, privateKey string, input []byte
    public static native String commitCallTx(String contract,String privateKey,byte[] input);

    public static native void initClient();
}
