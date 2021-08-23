# Signature library 
> support android, ios, dart, java, go, golang

## 1. how to use in java
```java
public class Test {
    static {
        System.loadLibrary("share_signer_main.so");
    }

    public static void main(String[] args) {
        Signer.initClient();
        //Generate address 
        System.out.println(Signer.genKey());
        
        String hex = Signer.transfer("{ \"contractAddr\": \"0xdb834d1f5baf312424fe3003524e2f5a52bf15b2\",\"to\":\"0xdb834d1f5baf312424fe3003524e2f5a52bf15b2\", \"base64PrivateKey\": \"h83oyA0Rc1GeAnzHnrdXBV8G1OKbreGjfMLG50XXtkLDPyWwc6/4jipJ4QTD0M2ZuVrjumMNR+sZJuxAGaVxig==\",\"amount\": 10000 }");    
        System.out.println(hex);
           
    }
}
```