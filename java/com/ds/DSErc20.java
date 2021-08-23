package com.ds;
public class DSErc20{
    private DSErc20(){

    }

    private DSErc20(String contract){
        setContract(contract);
    }

    public static DSErc20 Of(String contract) {
        return new DSErc20(contract);
    }
    //set ds contract 
    private native void setContract(String contract);
    //show had stacked of account
    public  native String stackOf(String addr);

    public  native String owner();

    public  native String addStack(String mainPrivateKey,String account,String amount);

    public  native String getMinerAddr(String peerId);

    //mainPrivateKey string, account string, amount string
    public native String penalty(String mainPrivateKey,String account,String amount);

    //mainPrivateKey, sender, recipient, amount string
    public native String exchange(String mainPrivateKey,String sender,String recipient,String amount);
}