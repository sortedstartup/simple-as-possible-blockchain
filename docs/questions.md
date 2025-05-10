1. we are creating the genesis block, who created the initial number of coins ? when are coins in created ? 

- hedera : it created fixed number of Hbars
- find for bitcoin ? i think during mining, relate it to halving.
- what is the merkle root hash of this block
    - since transaction tree is null

special things about generic block
prev block -> null

scenario 
-> in the genesis transaction lets say we create a transction from null -> speicl account transferring 1 B coins
 -> now how do we put it into circulation
 - airdrop 

- How to bootstrap a network of nodes


2. After submit transcion you validate the transaction and add it to your memory pool ?
2.1 It is not added to any block yet, how do the miner get this transaction from the memory pool.
2.2 What is the same person submits the same transaction to a different.

...

Finality - when the transction is confirmed

3. in the submit transaction, how do i know is sender has sufficnent balance
4. in the submit transaction, is the account in the system exists or not

5. How to create account?
   --> bit coin paper --> generate public private key pair before every transaction for privacy.
-> bitcoin, ethereum, hbar

Find answers for bitcoin, etherum and hedera.

## Accounts, Address, public key
Q) How do we create account on the bitcoin blockchain
A) We dont, there is no need to create any account on bitcoin.
The system is driven by public - private key.
A address is derived from any public key
address = simplefn (public_key)

This means if you generate a public-private key on a personal laptop,
you have a bitcoin address !
you dont need to create it on bitcoin via an API.

Q) If address=func(public_key) on a personal computer, how will a sender transfer it to me, how will I able to use that amount as a sender ?
The sender only need to know you r public key to send you money.
they will calcualte the address = func (public_key) 
The bitcoin network will check if the address is a valid public key and thats about it,
it is ready to transfer it to that address.
This transaction is now recorded on chain.
i.e it got recorded that a public key --transferred--to--> another public key
thats all bitcoin cares about.

How to claim money associated with your public derived address.
Private key !
private key is like a password and more than a password,
without revealing the private key you can proove using algorithms that you own it.
So you simply sign your transactions using your private key, the bitcoin system will recognize you are the owner of the public-private key pair.
So use your private key to move funds assigned to your public key and it works !


