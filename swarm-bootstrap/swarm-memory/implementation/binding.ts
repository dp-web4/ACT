// Ed25519 Binding Implementation
import { ed25519 } from '@noble/curves/ed25519';

export class LCTBinding {
  static async createBinding(entityData: any) {
    const privateKey = ed25519.utils.randomPrivateKey();
    const publicKey = ed25519.getPublicKey(privateKey);
    const signature = await ed25519.sign(entityData, privateKey);
    return { publicKey, signature };
  }
}