import main.Work;

public class Main {
    public static void main(String[] args) {
        Work db = new Work();

        db.add("1", "Hello, world!", 1);

        System.out.println(db.viewAll());
    }
}