// 外观模式能为复杂系统、 程序库或框架提供一个简单 （但有限） 的接口.

/// Facade hides a complex logic behind the API.
struct WalletFacade {
    account: Account,
    wallet: Wallet,
}

impl WalletFacade {
    pub fn new(account_id: String) -> Self {
        println!("Starting create account");
        let res = Self {
            account: Account::new(account_id),
            wallet: Wallet::new(),
        };
        println!("Account created");
        res
    }

    pub fn add_money_to_wallet(&mut self, account_id: &String, amount: u32) -> Result<(), String> {
        println!("Starting add money to wallet");
        self.account.check(account_id)?;
        self.wallet.credit_balance(amount);
        println!("Money added to wallet");
        Ok(())
    }

    pub fn deduct_money_from_wallet(
        &mut self,
        account_id: &String,
        amount: u32,
    ) -> Result<(), String> {
        println!("Starting deduct money from wallet");
        self.account.check(account_id)?;
        self.wallet.debit_balance(amount);
        println!("Money deducted from wallet");
        Ok(())
    }
}

struct Account {
    name: String,
}

impl Account {
    pub fn new(name: String) -> Self {
        Account { name }
    }

    pub fn check(&self, name: &String) -> Result<(), String> {
        if &self.name != name {
            return Err("Account name is incorrect".into());
        }
        println!("Account verified");
        Ok(())
    }
}

struct Wallet {
    balance: u32,
}

impl Wallet {
    pub fn new() -> Self {
        Wallet { balance: 0 }
    }

    pub fn credit_balance(&mut self, amount: u32) {
        self.balance += amount;
    }

    pub fn debit_balance(&mut self, amount: u32) {
        self.balance
            .checked_sub(amount)
            .expect("Insufficient balance");
    }
}

fn main() -> Result<(), String> {
    let mut wallet = WalletFacade::new("my_account".into());
    println!();

    wallet.add_money_to_wallet(&"my_account".into(), 50)?;
    println!();

    wallet.deduct_money_from_wallet(&"my_account".into(), 10)?;

    Ok(())
}
