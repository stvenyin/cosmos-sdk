# Clients

::: tip Note
🚧 We are actively working on documentation for SDK clients.
:::

## Gaia CLI

::: tip Note
🚧 We are actively working on improving documentation for Gaiacli and Gaiad.
:::

`gaiacli` is the command line interface to manage accounts and transactions on Cosmos testnets. Here is a list of useful `gaiacli` commands, including usage examples.

### Keys

#### Key Types

There are three types of key representations that are used:

- `cosmosaccaddr`
  - Derived from account keys generated by `gaiacli keys add`
  - Used to receive funds
  - e.g. `cosmosaccaddr15h6vd5f0wqps26zjlwrc6chah08ryu4hzzdwhc`
- `cosmosaccpub`
  - Derived from account keys generated by `gaiacli keys add`
  - e.g. `cosmosaccpub1zcjduc3q7fu03jnlu2xpl75s2nkt7krm6grh4cc5aqth73v0zwmea25wj2hsqhlqzm`
- `cosmosvalpub`
  - Generated when the node is created with `gaiad init`.
  - Get this value with `gaiad tendermint show-validator`
  - e.g. `cosmosvalpub1zcjduc3qcyj09qc03elte23zwshdx92jm6ce88fgc90rtqhjx8v0608qh5ssp0w94c`

#### Generate Keys

You'll need an account private and public key pair \(a.k.a. `sk, pk` respectively\) to be able to receive funds, send txs, bond tx, etc.

To generate a new key \(default _ed25519_ elliptic curve\):

```bash
gaiacli keys add <account_name>
```

Next, you will have to create a passphrase to protect the key on disk. The output of the above command will contain a _seed phrase_. Save the _seed phrase_ in a safe place in case you forget the password!

If you check your private keys, you'll now see `<account_name>`:

```bash
gaiacli keys show <account_name>
```

You can see all your available keys by typing:

```bash
gaiacli keys list
```

View the validator pubkey for your node by typing:

```bash
gaiad tendermint show-validator
```

::: danger Warning
We strongly recommend _NOT_ using the same passphrase for multiple keys. The Tendermint team and the Interchain Foundation will not be responsible for the loss of funds.
:::

### Account

#### Get Tokens

The best way to get tokens is from the [Cosmos Testnet Faucet](https://faucetcosmos.network). If the faucet is not working for you, try asking [#cosmos-validators](https://riot.im/app/#/room/#cosmos-validators:matrix.org). The faucet needs the `cosmosaccaddr` from the account you wish to use for staking.

#### Query Account balance

After receiving tokens to your address, you can view your account's balance by typing:

```bash
gaiacli account <account_cosmosaccaddr>
```

::: warning Note
When you query an account balance with zero tokens, you will get this error: `No account with address <account_cosmosaccaddr> was found in the state.` This can also happen if you fund the account before your node has fully synced with the chain. These are both normal.

:::

### Send Tokens

```bash
gaiacli send \
  --amount=10faucetToken \
  --chain-id=<chain_id> \
  --name=<key_name> \
  --to=<destination_cosmosaccaddr>
```

::: warning Note
The `--amount` flag accepts the format `--amount=<value|coin_name>`.
:::

Now, view the updated balances of the origin and destination accounts:

```bash
gaiacli account <account_cosmosaccaddr>
gaiacli account <destination_cosmosaccaddr>
```

You can also check your balance at a given block by using the `--block` flag:

```bash
gaiacli account <account_cosmosaccaddr> --block=<block_height>
```

### Staking

#### Set up a Validator

Please refer to the [Validator Setup](https://cosmos.network/docs/validators/validator-setup.html) section for a more complete guide on how to set up a validator-candidate.

#### Delegate to a Validator

On the upcoming mainnet, you can delegate `atom` to a validator. These [delegators](/resources/delegators-faq) can receive part of the validator's fee revenue. Read more about the [Cosmos Token Model](https://github.com/cosmos/cosmos/raw/master/Cosmos_Token_Model.pdf).

##### Query Validators

You can query the list of all validators of a specific chain:

```bash
gaiacli stake validators
```

If you want to get the information of a single validator you can check it with:

```bash
gaiacli stake validator <account_cosmosaccaddr>
```

#### Bond Tokens

On the testnet, we delegate `steak` instead of `atom`. Here's how you can bond tokens to a testnet validator (*i.e.* delegate):

```bash
gaiacli stake delegate \
  --amount=10steak \
  --validator=$(gaiad tendermint show-validator) \
  --name=<key_name> \
  --chain-id=<chain_id>
```

While tokens are bonded, they are pooled with all the other bonded tokens in the network. Validators and delegators obtain a percentage of shares that equal their stake in this pool.

::: tip Note
Don't use more `steak` thank you have! You can always get more by using the [Faucet](https://faucetcosmos.network/)!
:::

##### Query Delegations

Once submitted a delegation to a validator, you can see it's information by using the following command:

```bash
gaiacli stake delegation \
	--address-delegator=<account_cosmosaccaddr> \
	--address-validator=$(gaiad tendermint show-validator)
```

Or if you want to check all your current delegations with disctinct validators:

```bash
gaiacli stake delegations <account_cosmosaccaddr>
```

You can also get previous delegation(s) status by adding the `--height` flag.

#### Unbond Tokens

If for any reason the validator misbehaves, or you just want to unbond a certain amount of tokens, use this following command. You can unbond a specific `shares-amount` (eg:`12.1`\) or a `shares-percent` (eg:`25`) with the corresponding flags.

```bash
gaiacli stake unbond begin \
  --address-validator=$(gaiad tendermint show-validator) \
  --shares-percent=100 \
  --from=<key_name> \
  --chain-id=<chain_id>
```

Later you must complete the unbonding process by using the `gaiacli stake unbond complete` command:

```bash
gaiacli stake unbond complete \
  --address-validator=$(gaiad tendermint show-validator) \
  --from=<key_name> \
  --chain-id=<chain_id>
```

##### Query Unbonding-Delegations

Once you begin an unbonding-delegation, you can see it's information by using the following command:

```bash
gaiacli stake unbonding-delegation \
	--address-delegator=<account_cosmosaccaddr> \
	--address-validator=$(gaiad tendermint show-validator) \
```

Or if you want to check all your current unbonding-delegations with disctinct validators:

```bash
gaiacli stake unbonding-delegations <account_cosmosaccaddr>
```

You can also get previous unbonding-delegation(s) status by adding the `--height` flag.

#### Redelegate Tokens

A redelegation is a type delegation that allows you to bond illiquid tokens from one validator to another:

```bash
gaiacli stake redelegate begin \
  --address-validator-source=$(gaiad tendermint show-validator) \
  --address-validator-dest=<account_cosmosaccaddr> \
  --shares-percent=50 \
  --from=<key_name> \
  --chain-id=<chain_id>
```

Here you can also redelegate a specific `shares-amount` or a  `shares-percent` with the corresponding flags.

Later you must complete the redelegation process by using the `gaiacli stake redelegate complete` command:

```bash
gaiacli stake unbond complete \
  --address-validator=$(gaiad tendermint show-validator) \
  --from=<key_name> \
  --chain-id=<chain_id>
```

##### Query Redelegations

Once you begin an redelegation, you can see it's information by using the following command:

```bash
gaiacli stake redelegation \
	--address-delegator=<account_cosmosaccaddr> \
	--address-validator-source=$(gaiad tendermint show-validator) \
	--address-validator-dest=<account_cosmosaccaddr> \
```

Or if you want to check all your current unbonding-delegations with disctinct validators:

```bash
gaiacli stake redelegations <account_cosmosaccaddr>
```

You can also get previous redelegation(s) status by adding the `--height` flag.

### Governance

Governance is the process from which users in the Cosmos Hub can come to consensus on software upgrades, parameters of the mainnet or on custom text proposals. This is done through voting on proposals, which will be submitted by `Atom` holders on the mainnet.

Some considerations about the voting process:

- Voting is done by bonded `Atom` holders on a 1 bonded `Atom` 1 vote basis
- Delegators inherit the vote of their validator if they don't vote
- **Validators MUST vote on every proposal**. If a validator does not vote on a proposal, they will be **partially slashed**
- Votes are tallied at the end of the voting period (2 weeks on mainnet). Each address can vote multiple times to update its `Option` value (paying the transaction fee each time), only the last casted vote will count as valid
- Voters can choose between options `Yes`, `No`, `NoWithVeto` and `Abstain`
  At the end of the voting period, a proposal is accepted if `(YesVotes/(YesVotes+NoVotes+NoWithVetoVotes))>1/2` and `(NoWithVetoVotes/(YesVotes+NoVotes+NoWithVetoVotes))<1/3`. It is rejected otherwise

For more information about the governance process and how it works, please check out the Governance module [specification](https://github.com/cosmos/cosmos-sdk/tree/develop/docs/spec/governance).

#### Create a Governance proposal

In order to create a governance proposal, you must submit an initial deposit along with the proposal details:

- `title`: Title of the proposal
- `description`: Description of the proposal
- `type`: Type of proposal. Must be of value _Text_ (types _SoftwareUpgrade_ and _ParameterChange_ not supported yet).

```bash
gaiacli gov submit-proposal \
  --title=<title> \
  --description=<description> \
  --type=<Text/ParameterChange/SoftwareUpgrade> \
  --proposer=<account_cosmosaccaddr> \
  --deposit=<40steak> \
  --from=<name> \
  --chain-id=<chain_id>
```

##### Query proposals

Once created, you can now query information of the proposal:

```bash
gaiacli gov query-proposal \
  --proposal-id=<proposal_id>
```

Or query all available proposals:

```bash
gaiacli gov query-proposals
```

You can also query proposals filtered by `voter` or `depositer` by using the corresponding flags.

#### Increase deposit

In order for a proposal to be broadcasted to the network, the amount deposited must be above a `minDeposit` value (default: `10 steak`). If the proposal you previously created didn't meet this requirement, you can still increase the total amount deposited to activate it. Once the minimum deposit is reached, the proposal enters voting period:

```bash
gaiacli gov deposit \
  --proposal-id=<proposal_id> \
  --depositer=<account_cosmosaccaddr> \
  --deposit=<200steak> \
  --from=<name> \
  --chain-id=<chain_id>
```

> _NOTE_: Proposals that don't meet this requirement will be deleted after `MaxDepositPeriod` is reached.

#### Vote on a proposal

After a proposal's deposit reaches the `MinDeposit` value, the voting period opens. Bonded `Atom` holders can then cast vote on it:

```bash
gaiacli gov vote \
  --proposal-id=<proposal_id> \
  --voter=<account_cosmosaccaddr> \
  --option=<Yes/No/NoWithVeto/Abstain> \
  --from=<name> \
  --chain-id=<chain_id>
```

##### Query vote

Check the vote with the option you just submitted:

```bash
gaiacli gov query-vote \
  --proposal-id=<proposal_id> \
  --voter=<account_cosmosaccaddr>
```

## Gaia-Lite

::: tip Note
🚧 We are actively working on documentation for Gaia-lite.
:::
