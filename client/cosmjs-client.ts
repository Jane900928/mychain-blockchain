// MyChain CosmJS Client
import { SigningStargateClient, StargateClient } from "@cosmjs/stargate";
import { DirectSecp256k1HdWallet } from "@cosmjs/proto-signing";
import { coins, GasPrice } from "@cosmjs/amino";

export class MyChainClient {
    private stargateClient: StargateClient | null = null;
    private signingClient: SigningStargateClient | null = null;
    private wallet: DirectSecp256k1HdWallet | null = null;
    private readonly rpcEndpoint: string;
    private readonly chainId: string;

    constructor(rpcEndpoint: string = "http://localhost:26657", chainId: string = "mychain-1") {
        this.rpcEndpoint = rpcEndpoint;
        this.chainId = chainId;
    }

    async connect(): Promise<void> {
        try {
            this.stargateClient = await StargateClient.connect(this.rpcEndpoint);
            console.log("Connected to MyChain successfully");
        } catch (error) {
            console.error("Failed to connect to MyChain:", error);
            throw error;
        }
    }

    async connectWithMnemonic(mnemonic: string): Promise<string> {
        try {
            this.wallet = await DirectSecp256k1HdWallet.fromMnemonic(mnemonic, {
                prefix: "mychain",
            });

            const [firstAccount] = await this.wallet.getAccounts();
            
            this.signingClient = await SigningStargateClient.connectWithSigner(
                this.rpcEndpoint,
                this.wallet,
                {
                    gasPrice: GasPrice.fromString("0.1umychain"),
                }
            );

            console.log("Wallet connected:", firstAccount.address);
            return firstAccount.address;
        } catch (error) {
            console.error("Failed to connect wallet:", error);
            throw error;
        }
    }

    async createUser(name: string, email: string, senderAddress: string): Promise<string> {
        if (!this.signingClient) {
            throw new Error("Signing client not connected");
        }

        const msg = {
            typeUrl: "/mychain.mychain.MsgCreateUser",
            value: {
                creator: senderAddress,
                name: name,
                email: email,
            },
        };

        const result = await this.signingClient.signAndBroadcast(
            senderAddress,
            [msg],
            "auto",
            "Create new user"
        );

        if (result.code !== 0) {
            throw new Error(`Transaction failed: ${result.rawLog}`);
        }

        return result.transactionHash;
    }

    async transferTokens(
        senderAddress: string,
        receiverAddress: string,
        amount: string,
        denom: string = "mychain"
    ): Promise<string> {
        if (!this.signingClient) {
            throw new Error("Signing client not connected");
        }

        const msg = {
            typeUrl: "/mychain.mychain.MsgTransferTokens",
            value: {
                sender: senderAddress,
                receiver: receiverAddress,
                amount: coins(amount, denom),
            },
        };

        const result = await this.signingClient.signAndBroadcast(
            senderAddress,
            [msg],
            "auto",
            "Transfer tokens"
        );

        if (result.code !== 0) {
            throw new Error(`Transaction failed: ${result.rawLog}`);
        }

        return result.transactionHash;
    }

    async mintTokens(
        minterAddress: string,
        amount: string,
        denom: string = "mychain"
    ): Promise<string> {
        if (!this.signingClient) {
            throw new Error("Signing client not connected");
        }

        const msg = {
            typeUrl: "/mychain.mychain.MsgMintTokens",
            value: {
                minter: minterAddress,
                amount: coins(amount, denom),
            },
        };

        const result = await this.signingClient.signAndBroadcast(
            minterAddress,
            [msg],
            "auto",
            "Mint tokens"
        );

        if (result.code !== 0) {
            throw new Error(`Transaction failed: ${result.rawLog}`);
        }

        return result.transactionHash;
    }

    async registerMiner(
        minerAddress: string,
        description: string,
        commission: string
    ): Promise<string> {
        if (!this.signingClient) {
            throw new Error("Signing client not connected");
        }

        const msg = {
            typeUrl: "/mychain.mychain.MsgRegisterMiner",
            value: {
                miner: minerAddress,
                description: description,
                commission: commission,
            },
        };

        const result = await this.signingClient.signAndBroadcast(
            minerAddress,
            [msg],
            "auto",
            "Register as miner"
        );

        if (result.code !== 0) {
            throw new Error(`Transaction failed: ${result.rawLog}`);
        }

        return result.transactionHash;
    }

    async getCurrentHeight(): Promise<number> {
        if (!this.stargateClient) {
            throw new Error("Client not connected");
        }

        return await this.stargateClient.getHeight();
    }

    async getBalance(address: string, denom: string = "mychain"): Promise<string> {
        if (!this.stargateClient) {
            throw new Error("Client not connected");
        }

        const balance = await this.stargateClient.getBalance(address, denom);
        return balance.amount;
    }

    disconnect(): void {
        if (this.stargateClient) {
            this.stargateClient.disconnect();
            this.stargateClient = null;
        }
        if (this.signingClient) {
            this.signingClient.disconnect();
            this.signingClient = null;
        }
        this.wallet = null;
        console.log("Disconnected from MyChain");
    }

    isConnected(): boolean {
        return this.stargateClient !== null;
    }
}

export class WalletUtils {
    static async generateMnemonic(): Promise<string> {
        const wallet = await DirectSecp256k1HdWallet.generate(12, {
            prefix: "mychain",
        });
        return wallet.mnemonic;
    }

    static isValidAddress(address: string): boolean {
        return address.startsWith("mychain1") && address.length === 45;
    }
}

export default MyChainClient;
