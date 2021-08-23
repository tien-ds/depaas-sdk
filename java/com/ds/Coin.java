package com.ds;

public class Coin {

    public static native String balance(String account);

    public static native String transfer(String base64PrivateKey, String to, String amount);

    public static native String mint(String base64PrivateKey, String to, String amount);

}