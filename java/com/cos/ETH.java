package com.cos;

public class ETH {

    // GetTokenBalance(url string, contractAddr string, addr string)
    public static native String getTokenBalance(String url, String contractAddr, String addr);

    // getEthBalance
    public static native String getEthBalance(String url, String addr);

    // getContractTokenInfo
    public static native String getContractTokenInfo(String url, String contractAddr);

    public static native String allowance(String url, String contractAddr,String addr);

    public static native String swapApprove(String contractAddr,String mainNetPrivateKey);
}