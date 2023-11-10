
using System;
using System.Windows.Forms;

public class GuessingGameForm : Form
{
    private int randomNumber;
    private int guessCount;
    private TextBox guessInput;
    private Button guessButton;
    private Label messageLabel;

    public GuessingGameForm()
    {
        this.randomNumber = new Random().Next(1, 101);
        this.guessCount = 0;
        InitializeComponents();
    }

    private void InitializeComponents()
    {
        this.guessInput = new TextBox
        {
            Size = new System.Drawing.Size(200, 23),
            Location = new System.Drawing.Point(10, 10)
        };

        this.guessButton = new Button
        {
            Text = "Guess",
            Size = new System.Drawing.Size(100, 23),
            Location = new System.Drawing.Point(220, 10)
        };
        this.guessButton.Click += new EventHandler(GuessButton_Click);

        this.messageLabel = new Label
        {
            Size = new System.Drawing.Size(320, 23),
            Location = new System.Drawing.Point(10, 40)
        };

        this.Controls.Add(this.guessInput);
        this.Controls.Add(this.guessButton);
        this.Controls.Add(this.messageLabel);

        this.Text = "Random Number Guessing Game";
        this.Size = new System.Drawing.Size(350, 120);
    }

    private void GuessButton_Click(object sender, EventArgs e)
    {
        guessCount++;
        int userGuess;
        bool isNumeric = int.TryParse(this.guessInput.Text, out userGuess);

        if (!isNumeric)
        {
            MessageBox.Show("Please enter a valid number.");
            return;
        }

        if (userGuess > randomNumber)
        {
            this.messageLabel.Text = "Too high, try again.";
        }
        else if (userGuess < randomNumber)
        {
            this.messageLabel.Text = "Too low, try again.";
        }
        else
        {
            this.messageLabel.Text = $"Congratulations! You guessed the number in {guessCount} attempts.";
            this.randomNumber = new Random().Next(1, 101);
            this.guessCount = 0;
        }

        this.guessInput.Text = "";
    }

    static void Main()
    {
        Application.Run(new GuessingGameForm());
    }
}
