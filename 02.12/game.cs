using System;
using System.Windows.Forms;
using System.Drawing;

public class RockPaperScissorsGame : Form
{
    private Button rockButton;
    private Button paperButton;
    private Button scissorsButton;
    private Label computerChoiceLabel;
    private Label resultLabel;

    private Random random = new Random();

    public RockPaperScissorsGame()
    {
        this.Text = "Rock, Paper, Scissors Game";
        this.Size = new Size(300, 200);

        rockButton = new Button { Text = "Rock", Location = new Point(10, 10) };
        paperButton = new Button { Text = "Paper", Location = new Point(10, 40) };
        scissorsButton = new Button { Text = "Scissors", Location = new Point(10, 70) };

        computerChoiceLabel = new Label { Text = "Computer's choice: ", Location = new Point(10, 100), AutoSize = true };
        resultLabel = new Label { Text = "Result: ", Location = new Point(10, 130), AutoSize = true };

        rockButton.Click += OnUserChoice;
        paperButton.Click += OnUserChoice;
        scissorsButton.Click += OnUserChoice;

        this.Controls.Add(rockButton);
        this.Controls.Add(paperButton);
        this.Controls.Add(scissorsButton);
        this.Controls.Add(computerChoiceLabel);
        this.Controls.Add(resultLabel);
    }

    private void OnUserChoice(object sender, EventArgs e)
    {
        Button clickedButton = (Button)sender;
        string userChoice = clickedButton.Text;
        string computerChoice = GetComputerChoice();
        string result = DetermineWinner(userChoice, computerChoice);
        computerChoiceLabel.Text = $"Computer's choice: {computerChoice}";
        resultLabel.Text = $"Result: {result}";
    }

    private string GetComputerChoice()
    {
        int choice = random.Next(1, 4);
        switch (choice)
        {
            case 1: return "Rock";
            case 2: return "Paper";
            case 3: return "Scissors";
            default: return "Rock";
        }
    }

    private string DetermineWinner(string userChoice, string computerChoice)
    {
        if (userChoice == computerChoice)
            return "It's a tie!";

        if ((userChoice == "Rock" && computerChoice == "Scissors") ||
            (userChoice == "Scissors" && computerChoice == "Paper") ||
            (userChoice == "Paper" && computerChoice == "Rock"))
        {
            return "You win!";
        }
        else
        {
            return "You lose!";
        }
    }

    public static void Main()
    {
        Application.EnableVisualStyles();
        Application.Run(new RockPaperScissorsGame());
    }
}
