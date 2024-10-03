package test;
import static org.junit.jupiter.api.Assertions.*;
import main.Work;

class WorkTest {
    Work db = new Work();
    @org.junit.jupiter.api.Test
    void add() {
        db.add("1", "Hello, world", 1);
        assertEquals(1, db.size());
    }

    @org.junit.jupiter.api.Test
    void delete() {
        db.add("1", "Hello, world", 1);
        db.delete("1");
        assertEquals(0, db.size());
    }

    @org.junit.jupiter.api.Test
    void put() {
        db.add("1", "Hello, world", 1);
        db.put("1", "World, Hello!");
        assertEquals("World, Hello!", db.get("1"));
    }

    @org.junit.jupiter.api.Test
    void viewAll() {
    }

    @org.junit.jupiter.api.Test
    void viewKey() {
        db.add("1", "Hello, world", 1);
        assertEquals("Hello, world", db.get("1"));
    }
}