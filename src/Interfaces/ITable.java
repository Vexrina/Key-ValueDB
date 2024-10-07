package Interfaces;

public interface ITable {
    void add(String key, String value, int time);
    void delete(String key);
    void put(String key, String value);
    String get(String key);
    int size();
}
