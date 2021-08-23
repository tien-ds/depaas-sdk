package com.cos;

public class Oracle {
    private String newContracts;

    public static native void gatewayInit();

    public Oracle(String sMainNetCoin, String sLoomCoin) {
        newContracts = newContracts(sMainNetCoin, sLoomCoin);
    }

    private native String newContracts(String sMainNetCoin, String sLoomCoin);

    
    private native String deposit(String ptr, String mainNetPrivateKey, String tokenAmount);

    
    private native String withdrawal(String ptr, String base64PrivateKey, String mainNetPrivateKey, String tokenAmount);

    // cos => eth
    public String withdrawal(String base64PrivateKey, String mainNetPrivateKey, String tokenAmount) {
        return withdrawal(newContracts, base64PrivateKey, mainNetPrivateKey, tokenAmount);
    }

    // eth => cos
    public String deposit(String mainNetPrivateKey, String tokenAmount){
        return deposit(newContracts,mainNetPrivateKey, tokenAmount);
    }
}