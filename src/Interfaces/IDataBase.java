package Interfaces;

public interface IDataBase {
    void add(String key, ITable value);
    ITable get(String key);
    void remove(String key);
    void put(String key, ITable value);
}
