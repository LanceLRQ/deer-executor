use std::io;
fn main() -> io::Result<()> {
    loop {
        let mut buffer = String::new();
        let bytes = io::stdin().read_line(&mut buffer)?;
        if bytes == 0 {
            break;
        }
        let nums: Vec<&str> = buffer.trim().split(' ').collect();
        let num: isize = nums[0].to_string().parse::<isize>().unwrap()
            + nums[1].to_string().parse::<isize>().unwrap();
        println!("{}", num);
    }
    Ok(())
}
