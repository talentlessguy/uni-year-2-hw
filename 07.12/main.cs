using System;
using System.IO;
using System.Windows.Forms;
using System.Drawing;

namespace DriverLicenseExamGrader
{
    public class MainForm : Form
    {
        private const int TotalQuestions = 20;
        private const int PassingScore = 15;
        private readonly char[] correctAnswers = { 'B', 'D', 'A', 'A', 'C', 'A', 'B', 'A', 'C', 'D', 'B', 'C', 'D', 'A', 'D', 'C', 'C', 'B', 'D', 'A' };

        private Button gradeButton;
        private Label resultLabel;
        private DataGridView dataGridView1;

        public MainForm()
        {
            InitializeComponents();
        }

        private void InitializeComponents()
        {
            this.Text = "Driver's License Exam Grader";
            this.Size = new Size(600, 400);

            gradeButton = new Button();
            gradeButton.Text = "Grade Exam";
            gradeButton.Location = new Point(20, 20);
            gradeButton.Click += GradeButton_Click;
            this.Controls.Add(gradeButton);

            resultLabel = new Label();
            resultLabel.Location = new Point(20, 60);
            this.Controls.Add(resultLabel);

            dataGridView1 = new DataGridView();
            dataGridView1.Location = new Point(20, 180);
            dataGridView1.Size = new Size(550, 150);
            dataGridView1.AllowUserToAddRows = false;
            dataGridView1.AllowUserToDeleteRows = false;
            dataGridView1.AllowUserToOrderColumns = true;
            dataGridView1.ReadOnly = true;

            DataGridViewTextBoxColumn questionNumberColumn = new DataGridViewTextBoxColumn
            {
                HeaderText = "Question Number",
                Name = "QuestionNumber"
            };
            DataGridViewTextBoxColumn correctAnswerColumn = new DataGridViewTextBoxColumn
            {
                HeaderText = "Correct Answer",
                Name = "CorrectAnswer"
            };
            DataGridViewTextBoxColumn studentAnswerColumn = new DataGridViewTextBoxColumn
            {
                HeaderText = "Student Answer",
                Name = "StudentAnswer"
            };

            dataGridView1.Columns.AddRange(new DataGridViewColumn[] { questionNumberColumn, correctAnswerColumn, studentAnswerColumn });

            this.Controls.Add(dataGridView1);
        }

        private void GradeButton_Click(object sender, EventArgs e)
        {
            OpenFileDialog openFileDialog = new OpenFileDialog();
            openFileDialog.Filter = "Text Files|*.txt";

            if (openFileDialog.ShowDialog() == DialogResult.OK)
            {
                string filePath = openFileDialog.FileName;
                string[] studentAnswers = File.ReadAllLines(filePath);

                if (studentAnswers.Length != TotalQuestions)
                {
                    MessageBox.Show("The file should contain answers for all 20 questions.", "Invalid File", MessageBoxButtons.OK, MessageBoxIcon.Error);
                    return;
                }

                int correctCount = 0;
                dataGridView1.Rows.Clear();

                for (int i = 0; i < TotalQuestions; i++)
                {
                    char correctAnswer = correctAnswers[i];
                    char studentAnswer = studentAnswers[i].ToUpper()[0];

                    if (correctAnswer == studentAnswer)
                    {
                        correctCount++;
                    }
                    else
                    {
                        dataGridView1.Rows.Add(i + 1, correctAnswer, studentAnswer);
                    }
                }

                if (correctCount >= PassingScore)
                {
                    resultLabel.Text = "Pass";
                }
                else
                {
                    resultLabel.Text = "Fail";
                }
            }
        }

        [STAThread]
        public static void Main()
        {
            Application.EnableVisualStyles();
            Application.SetCompatibleTextRenderingDefault(false);
            Application.Run(new MainForm());
        }
    }
}
