import os

ADMIN_EMAIL = "admin@elimt.com"
ADMIN_PASSWORD = "strongpassword"


def autoGenerateCode(length):
    import random
    import string

    return "".join(
        random.choice(string.ascii_uppercase + string.ascii_lowercase + string.digits)
        for _ in range(length)
    )


def save_env(email, password):
    with open(".env", "w") as f:
        f.write(f"ADMIN_EMAIL={email}\nADMIN_PASSWORD={password}")


def email_validate(email):
    import re

    regex = r"\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Z|a-z]{2,}\b"
    return re.fullmatch(regex, email) is not None


def email_input():
    print("\n[Step 1] Please enter the admin email address:")
    print("(Leave this empty to generate a fake email address)")

    email = input("Your email: ").strip()

    if not email:
        email = f"{autoGenerateCode(6)}@elimt.com"
        print(f"Generated email: {email}")
    elif not email_validate(email):
        print("Invalid email address! Please try again.")
        return email_input()

    return email


def password_input():
    print("\n[Step 2] Please enter the admin password (at least 8 characters):")
    print("(Leave this empty to generate a secure password)")

    password = input("Your password: ").strip()

    if not password:
        password = autoGenerateCode(16)
        print(f"Generated password: {password}")
    elif len(password) < 8:
        print("Password must be at least 8 characters! Please try again.")
        return password_input()

    return password


def confirm_input(email, password):
    print("\n[Step 3] Confirm your information:")
    print("=" * 60)
    print(f"Email: {email}")
    print(f"Password: {password}")
    print("=" * 60)

    print("\n1. Confirm and save")
    print("2. Re-enter email and password")
    print("3. Exit setup")

    choice = input("Select an option (1/2/3): ").strip()

    if choice == "1":
        return True
    elif choice == "2":
        return False
    elif choice == "3":
        print("Exiting setup. No changes made.")
        exit()
    else:
        print("Invalid input! Please try again.")
        return confirm_input(email, password)


def main():
    os.system("cls" if os.name == "nt" else "clear")
    print("Welcome to the ELIMT setup script!")

    while True:
        email = email_input()
        password = password_input()

        if confirm_input(email, password):
            save_env(email, password)
            os.system("cls" if os.name == "nt" else "clear")
            print("\nSetup completed successfully!")
            print("You can now run the server with the following command:\n")
            print("docker-compose up -d\n")
            break


if __name__ == "__main__":
    main()
